package streamserver

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"net/http"
	"os"
	"time"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/a13labs/m3uproxy/pkg/xmltv/types"
	"github.com/gorilla/mux"

	"github.com/a13labs/m3uproxy/pkg/xmltv"
)

type EPGHandler struct {
	config      *ServerConfig
	sources     []string
	lastUpdated time.Time
}

func NewEPGHandler(config *ServerConfig) *EPGHandler {
	return &EPGHandler{
		config:  config,
		sources: config.GetEpg(),
	}
}

func (e *EPGHandler) RegisterRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/epg/{channel_id}", basicAuth(e.channelEpgRequest))
	return r
}

func (e *EPGHandler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(e.config.data.EPGScanTime) * time.Second)
		defer ticker.Stop()
		e.UpdateSources(ctx)
		for {
			select {
			case <-ticker.C:
				logger.Infof("Starting EPG sources update")
				e.UpdateSources(ctx)
			case <-ctx.Done():
				logger.Info("EPG update routine stopping due to context cancellation")
				return
			}
		}
	}()
}

func (e *EPGHandler) UpdateSources(ctx context.Context) {

	newUpdateTTime := time.Now().Add(time.Duration(e.config.data.EPGScanTime) * time.Second)
	if e.lastUpdated.After(newUpdateTTime) {
		logger.Debugf("Skipping EPG update: last update was at %v, next scheduled for %v", e.lastUpdated, newUpdateTTime)
		return
	}

	for _, src := range e.sources {
		select {
		case <-ctx.Done():
			logger.Info("EPG update routine stopping due to context cancellation")
			return
		default:
			if src == "" {
				logger.Warn("Skipping empty EPG source")
				continue
			}

			content, err := loadContent(src)
			if err != nil {
				logger.Errorf("EPG file not found at %s: %v", src, err)
				continue
			}

			var epg xmltv.EPG
			if err := xml.Unmarshal([]byte(content), &epg); err != nil {
				logger.Errorf("Failed to parse XMLTV from %s: %v", src, err)
				continue
			}

			logger.Infof("Loaded EPG from %s: %d channels, %d programmes", src, len(epg.Channels), len(epg.Programmes))

			for _, channel := range epg.Channels {
				select {
				case <-ctx.Done():
					logger.Info("EPG update routine stopping due to context cancellation")
					return
				default:
					channelEPG := xmltv.EPG{
						Channels:   []xmltv.Channel{channel},
						Programmes: []xmltv.Programme{},
					}

					for _, programme := range epg.Programmes {
						if programme.Channel == channel.ID {
							channelEPG.Programmes = append(channelEPG.Programmes, programme)
						}
					}

					channelXML, err := xml.MarshalIndent(channelEPG, "", "  ")
					if err != nil {
						logger.Errorf("Error generating channel EPG for channel %s: %v", channel.ID, err)
						continue
					}

					filename := e.config.data.CacheDir + "/" + base64.StdEncoding.EncodeToString([]byte(channel.ID)) + newUpdateTTime.Format("_20060102150405") + "_epg.xml"
					if err := os.WriteFile(filename, channelXML, 0644); err != nil {
						logger.Errorf("Failed to write EPG file for channel %s to %s: %v", channel.ID, filename, err)
						continue
					}
					logger.Debugf("Saved EPG for channel %s to %s", channel.ID, filename)
				}
			}
		}
	}
	e.lastUpdated = newUpdateTTime
	logger.Infof("EPG sources updated at %v", e.lastUpdated)

	files, err := os.ReadDir(e.config.data.CacheDir)
	if err != nil {
		logger.Errorf("Error reading cache directory %s: %v", e.config.data.CacheDir, err)
		return
	}

	for _, file := range files {
		select {
		case <-ctx.Done():
			logger.Debug("EPG cleanup routine stopping due to context cancellation")
			return
		default:
			if file.IsDir() {
				continue
			}
			info, err := file.Info()
			if err != nil {
				logger.Debugf("Error getting info for file %s: %v", file.Name(), err)
				continue
			}
			// Remove files older than the current newUpdateTTime
			// (to avoid removing files just created in this update cycle)
			//
			if time.Since(info.ModTime()) > time.Duration(e.config.data.EPGScanTime)*time.Second && info.ModTime().Before(newUpdateTTime) {
				err := os.Remove(e.config.data.CacheDir + "/" + file.Name())
				if err != nil {
					logger.Debugf("Failed to remove old EPG file %s: %v", file.Name(), err)
				} else {
					logger.Debugf("Removed old EPG file %s", file.Name())
				}
			}
		}
	}
}

func (e *EPGHandler) channelEpgRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Warnf("Rejected %s request to /epg: only GET is allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	channelID := vars["channel_id"]
	timestamp := r.URL.Query().Get("timestamp")
	var reqTimeRaw *time.Time
	if timestamp != "" {
		t, err := parseTimestamp(timestamp)
		if err != nil {
			logger.Warnf("Invalid timestamp format in request: %s", timestamp)
			http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
			return
		}
		reqTimeRaw = t
	}

	channelFile := e.config.data.CacheDir + "/" + base64.StdEncoding.EncodeToString([]byte(channelID)) + e.lastUpdated.Format("_20060102150405") + "_epg.xml"

	if _, err := os.Stat(channelFile); os.IsNotExist(err) {
		logger.Warnf("EPG file for channel %s not found at %s", channelID, channelFile)
		http.Error(w, "EPG for channel not found", http.StatusNotFound)
		return
	}

	content, err := loadContent(channelFile)
	if err != nil {
		logger.Errorf("Failed to load EPG file at %s: %v", channelFile, err)
		http.Error(w, "EPG file not found", http.StatusNotFound)
		return
	}

	var epg xmltv.EPG
	if err := xml.Unmarshal([]byte(content), &epg); err != nil {
		logger.Errorf("Failed to parse XMLTV for channel %s: %v", channelID, err)
		http.Error(w, "Error parsing EPG", http.StatusInternalServerError)
		return
	}

	filteredEPG := xmltv.EPG{
		Channels:   epg.Channels,
		Programmes: []xmltv.Programme{},
	}

	for _, programme := range epg.Programmes {
		if programme.Channel != channelID {
			continue
		}

		if timestamp != "" && reqTimeRaw != nil {
			// Use programme Start/Stop times directly (avoid formatting + reparsing)
			start := programme.Start.Time
			if programme.Stop == nil {
				logger.Warnf("Skipping programme without stop time for channel %s", channelID)
				continue
			}
			stop := programme.Stop.Time

			// Compare request time in programme's timezone
			reqTime := *reqTimeRaw
			reqInProgTZ := reqTime.In(start.Location())
			if reqInProgTZ.Before(start) || reqInProgTZ.After(stop) {
				continue
			}

			// Convert programme times to request timezone for output (representation only)
			convertedStart := start.In(reqTime.Location())
			convertedStop := stop.In(reqTime.Location())
			programme.Start = types.XMLTVTime{Time: convertedStart}
			programme.Stop = &types.XMLTVTime{Time: convertedStop}
		}

		filteredEPG.Programmes = append(filteredEPG.Programmes, programme)
	}

	logger.Debugf("Serving EPG for channel %s: %d channels, %d programmes", channelID, len(epg.Channels), len(filteredEPG.Programmes))
	filteredXML, err := xml.MarshalIndent(filteredEPG, "", "  ")
	if err != nil {
		logger.Errorf("Failed to generate filtered EPG for channel %s: %v", channelID, err)
		http.Error(w, "Error generating filtered EPG", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(filteredXML))
}

func parseTimestamp(ts string) (*time.Time, error) {
	tt, err := time.Parse("20060102150405 -0700", ts)
	if err != nil {
		tt, err = time.Parse("200601021504", ts)
		if err != nil {
			return nil, err
		}
	}
	return &tt, nil
}
