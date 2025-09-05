package streamserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/a13labs/m3uproxy/pkg/auth"
	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/oschwald/geoip2-golang"

	"github.com/gorilla/mux"
)

type StreamServer struct {
	running            bool
	updateTimer        *time.Timer
	config             *ServerConfig
	api                *APIHandler
	epg                *EPGHandler
	channels           *ChannelsHandler
	player             *PlayerHandler
	restartChan        chan bool
	router             *mux.Router
	geoipDb            *geoip2.Reader
	geoipWhitelist     map[string]bool
	geoIPCidrWhitelist []*net.IPNet
}

// Initialize the server
func NewStreamServer(configPath string) *StreamServer {
	s := StreamServer{
		restartChan: make(chan bool),
		config:      NewServerConfig(configPath),
		router:      mux.NewRouter(),
	}

	return &s
}

func (s *StreamServer) Run() {

	// Initialize logger with log level support
	if s.config.data.LogLevel != "" {
		logger.InitWithLevel(s.config.data.LogFile, s.config.data.LogLevel)
	} else {
		logger.Init(s.config.data.LogFile)
	}

	s.channels = NewChannelsHandler(s.config)
	s.api = NewAPIHandler(s.config, &s.restartChan, s.channels)
	s.api.RegisterRoutes(s.router)
	s.channels.RegisterRoutes(s.router)

	s.epg = NewEPGHandler(s.config)
	s.epg.RegisterRoutes(s.router)

	s.player = NewPlayerHandler(s.config)
	s.player.RegisterRoutes(s.router)

	s.router.HandleFunc("/health", s.healthCheckRequest)

	for {
		logger.Infof("Starting M3U Proxy Server")

		logger.Infof("Starting stream server")
		logger.Infof("Playlist: %s", s.config.data.Playlist)
		logger.Infof("EPG: %s", s.config.data.Epg)

		err := auth.InitializeAuth(s.config.data.Auth)
		if err != nil {
			logger.Errorf("Failed to initialize authentication: %s", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		s.updateTimer = time.NewTimer(time.Duration(s.config.data.ScanTime) * time.Second)
		s.running = true
		go func() {
			s.channels.Load(ctx)
			for {
				<-s.updateTimer.C
				s.channels.Load(ctx)
				if s.running {
					s.updateTimer.Reset(time.Duration(s.config.data.ScanTime) * time.Second)
				}
			}
		}()

		if s.configureSecurity() != nil {
			logger.Warn("GeoIP database not found, geo-location will not be available.")
		}

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", s.config.data.Port),
			Handler: s.secure(s.router),
		}

		// Channel to listen for termination signal (SIGINT, SIGTERM)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		go func() {
			logger.Infof("Server listening on %s.", server.Addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Errorf("Server failed: %v", err)
			}
		}()

		quitServer := false
		select {
		case <-sigChan:
			logger.Info("Signal received, shutting down server...")
			quitServer = true
			cancel()
		case <-s.restartChan:
			quitServer = false
		}

		s.updateTimer.Stop()
		s.running = false
		logger.Info("Stream server stopped")

		s.cleanGeoIp()

		if err := server.Shutdown(ctx); err != nil {
			logger.Errorf("Server forced to shutdown: %v", err)
		}

		if quitServer {
			logger.Info("Server shutdown.")
			break
		}
	}
}

func (s *StreamServer) Restart() {
	go func() {
		logger.Info("Server restart in 3 seconds")
		time.Sleep(3 * time.Second)
		s.restartChan <- true
	}()
}

func (s *StreamServer) healthCheckRequest(w http.ResponseWriter, r *http.Request) {
	if s.channels == nil {
		http.Error(w, "Streams not loaded", http.StatusServiceUnavailable)
		return
	}
	// health := map[string]interface{}{
	//     "status":        "ok",
	//     "activeStreams": activeStreams,
	//     "uptime":        time.Since(startTime).String(),
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(health)

	w.WriteHeader(http.StatusOK)
}

func (s *StreamServer) configureSecurity() error {

	var err error

	if s.config.data.Security.GeoIP.Database == "" {
		return nil
	}

	s.geoipDb, err = geoip2.Open(s.config.data.Security.GeoIP.Database)
	if err != nil {
		s.geoipDb = nil
		return err
	}

	s.geoipWhitelist = make(map[string]bool)
	for _, country := range s.config.data.Security.GeoIP.Whitelist {
		s.geoipWhitelist[country] = true
	}

	s.geoIPCidrWhitelist = make([]*net.IPNet, 0)

	for _, cidr := range s.config.data.Security.GeoIP.InternalNetworks {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		s.geoIPCidrWhitelist = append(s.geoIPCidrWhitelist, ipnet)
	}

	return nil
}

func (s *StreamServer) cleanGeoIp() {
	if s.geoipDb != nil {
		s.geoipDb.Close()
	}
}

func (s *StreamServer) secure(next http.Handler) http.Handler {

	if len(s.config.data.Security.AllowedCORSDomains) > 0 {
		next = s.cors(next)
	}

	if s.geoipDb == nil {
		return next
	}

	logger.Info("GeoIP enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := ""
		if r.Header.Get("X-Real-IP") != "" {
			ip = r.Header.Get("X-Real-IP")
		} else if r.Header.Get("X-Forwarded-For") != "" {
			ips := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
			ip = ips[0]
		} else {
			var err error
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if ip == "" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		parsedIP := net.ParseIP(ip)

		for _, ipnet := range s.geoIPCidrWhitelist {
			if ipnet.Contains(parsedIP) {
				next.ServeHTTP(w, r)
				return
			}
		}

		record, err := s.geoipDb.Country(parsedIP)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		countryCode := record.Country.IsoCode
		if _, ok := s.geoipWhitelist[countryCode]; !ok {
			logger.Infof("Access Denied: %s, Country: %s", ip, countryCode)
			http.Error(w, "Access Denied", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *StreamServer) cors(next http.Handler) http.Handler {
	logger.Info("CORS enabled")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", strings.Join(s.config.data.Security.AllowedCORSDomains, ","))
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
