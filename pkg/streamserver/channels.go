package streamserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/a13labs/m3uproxy/pkg/auth"
	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	"github.com/a13labs/m3uproxy/pkg/provider"
	"github.com/a13labs/m3uproxy/pkg/sources"
	"github.com/gorilla/mux"
)

var (
	// licenseManger holds per-process DRM license keys discovered while
	// parsing playlist entries. It is intentionally package-level so the
	// handler workers can register licenses discovered during parsing.
	// Note: the variable name preserves historical spelling for
	// compatibility within this package.
	licenseManger *streamLicenseManager
)

// streamEntry groups metadata and sources for a single logical channel.
//
// Fields:
//   - index: original order in the playlist (used for stable sorting)
//   - tvgId: identifier used for EPG matching and playlist URIs
//   - sources: collection of backend sources for the channel
type streamEntry struct {
	tvgId   string
	sources sources.Sources
	entry   m3uparser.M3UEntry
}

// ChannelsHandler manages the in-memory set of channels and their sources
// and provides HTTP handlers to serve playlists, manifests and media.
type ChannelsHandler struct {
	config         *ServerConfig
	playlistConfig *provider.PlaylistConfig
	channelsMux    sync.RWMutex
	channels       []*streamEntry
	channelsIdMap  map[string]int
}

func NewChannelsHandler(config *ServerConfig) *ChannelsHandler {
	return &ChannelsHandler{
		config:        config,
		channels:      make([]*streamEntry, 0),
		channelsIdMap: make(map[string]int),
		channelsMux:   sync.RWMutex{},
	}
}

func (p *ChannelsHandler) RegisterRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/channels.m3u", basicAuth(p.playlistRequest))
	r.HandleFunc("/drm/licensing", basicAuth(licenseKeysRequest))
	r.HandleFunc("/{token}/{channelId}/media/{path:.*}", p.mediaRequest)
	r.HandleFunc("/{token}/{channelId}/{path:.*}", p.manifestRequest)
	return r
}

func (p *ChannelsHandler) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Duration(p.config.data.ScanTime) * time.Second)
		defer ticker.Stop()
		p.UpdateSources(ctx)
		for {
			select {
			case <-ticker.C:
				// Periodic update of playlist sources from configured
				// playlist provider(s).
				logger.Infof("EPG sources update triggered")
				p.UpdateSources(ctx)
			case <-ctx.Done():
				logger.Info("Channels update routine stopping: context canceled")
				return
			}
		}
	}()
}

func (p *ChannelsHandler) UpdateSources(ctx context.Context) error {
	var err error
	p.playlistConfig, err = provider.LoadPlaylistConfig(p.config.data.Playlist)
	if err != nil {
		return err
	}

	m3uCache, err := provider.Load(p.playlistConfig)
	if err != nil {
		return err
	}

	// Load licenses found in playlist entries.
	// Current implementation supports clearkey licenses via KODIPROP tags.
	for _, entry := range m3uCache.Entries {
		keyType, keyId, keyValue := "", "", ""
		for _, tag := range entry.Tags {
			if tag.Tag == "KODIPROP" {
				if strings.HasPrefix(tag.Value, "inputstream.adaptive.license_type=") {
					parts := strings.Split(tag.Value, "=")
					if len(parts) == 2 {
						keyType = parts[1]
					}
					continue
				}
				if strings.HasPrefix(tag.Value, "inputstream.adaptive.license_key=") {

					if keyType == "org.w3.clearkey" {
						parts := strings.Split(tag.Value, "=")
						if len(parts) == 2 {
							licenseKey := parts[1]
							keyId = strings.Split(licenseKey, ":")[0]
							keyValue = strings.Split(licenseKey, ":")[1]

							if licenseManger == nil {
								licenseManger = newStreamLicenseManager()
							}
							logger.Infof("Discovered CLEARKEY license: adding key id=%s", keyId)
							licenseManger.addLicense("clearkey", keyId, keyValue)
							keyType, keyId, keyValue = "", "", ""
							break
						}
					}
				}
			}
		}
	}

	logger.Infof("Loaded %d playlist entries from %s", m3uCache.StreamCount(), p.config.data.Playlist)

	var wg sync.WaitGroup
	streamsChan := make(chan *streamEntry)
	stopWorkers := make(chan bool)

	for i := 0; i < p.config.data.NumWorkers; i++ {
		wg.Add(1)
		go monitorWorker(streamsChan, stopWorkers, &wg)
	}

	go func() {
		availableChannels := make(map[string]*streamEntry)
		availableChannelsIdMap := make(map[string]int)
		channelBucket := make([]string, 0)
		for i, entry := range m3uCache.Entries {
			select {
			case <-ctx.Done():
				stopWorkers <- true
				wg.Wait()
				return
			default:
				if entry.URI == "" {
					// Skip entries without URIs
					continue
				}

				tvgId := entry.ExtInfTags.GetValue("tvg-id")
				if tvgId == "" {
					// Fall back to the title when no tvg-id is provided
					tvgId = entry.Title
				}

				radio := entry.ExtInfTags.GetValue("radio")
				if tvgId == "" && radio == "" {
					logger.Warnf("Skipping entry without tvg-id or radio: %s", entry.URI)
					continue
				}

				_, ok := availableChannels[tvgId]
				if !ok {
					// New channel
					availableChannels[tvgId] = &streamEntry{
						tvgId:   tvgId,
						sources: sources.NewSources(),
						entry:   entry,
					}
					availableChannelsIdMap[tvgId] = i
					channelBucket = append(channelBucket, tvgId)
				}
			}
		}

		if len(p.config.data.ChannelBucket) > 0 {
			logger.Infof("Channel bucket configured with %d channels, filtering available channels", len(p.config.data.ChannelBucket))
			channelBucket = p.config.data.ChannelBucket
		}

		for i, channelId := range channelBucket {
			select {
			case <-ctx.Done():
				stopWorkers <- true
				wg.Wait()
				return
			default:
				if channel, ok := availableChannels[channelId]; ok {

					id, ok := p.channelsIdMap[channelId]
					if ok {
						if id != i {
							// Switch positions
							p.channelsMux.Lock()
							p.channels[i], p.channels[id] = p.channels[id], p.channels[i]
							p.channels[i] = channel
							p.channelsIdMap[channelId] = i
							p.channelsMux.Unlock()
						} else {
							// Same position, do nothing
						}
					} else {
						// New channel, add it
						p.channelsMux.Lock()
						p.channels = append(p.channels, channel)
						p.channelsIdMap[channelId] = i
						p.channelsMux.Unlock()
					}

					channel := p.channels[i]

					if channel.sources.SourceExists(channel.entry) {
						logger.Warnf("Stream source already exists (skipping): %s", channel.entry.URI)
						continue
					}

					logger.Infof("Adding stream source uri=%s to channel=%s", channel.entry.URI, channel.tvgId)
					added, err := channel.sources.AddSource(channel.entry, p.config.data.Timeout)
					if err != nil {
						logger.Errorf("Failed to add stream source uri=%s: %v", channel.entry.URI, err)
						continue
					}
					if !added {
						logger.Warnf("Stream source was not added (duplicate): %s", channel.entry.URI)
						continue
					}

					streamsChan <- channel
				} else {
					logger.Warnf("Channel %s not found, skipping (not in playlist?)", channelId)
				}
			}
		}

		// Remove channels that are no longer in the bucket
		// We can remove everything above len(ChannelBucket)
		p.channelsMux.Lock()
		p.channels = p.channels[:len(channelBucket)]
		newIdMap := make(map[string]int)
		for i, channel := range channelBucket {
			newIdMap[channel] = i
		}
		p.channelsIdMap = newIdMap
		p.channelsMux.Unlock()
		close(streamsChan)
	}()

	wg.Wait()

	return nil
}

func (p *ChannelsHandler) playlistRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	authParts := strings.SplitN(authHeader, " ", 2)
	token := authParts[1]

	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = r.URL.Scheme
	}
	if scheme == "" {
		scheme = "http"
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("#EXTM3U\n"))

	// Write the playlist
	p.channelsMux.RLock()
	defer p.channelsMux.RUnlock()
	for _, channel := range p.channels {
		if !channel.sources.Active() {
			continue
		}

		tvgId := strings.ReplaceAll(channel.tvgId, " ", "%20")
		uri := fmt.Sprintf("%s://%s/%s/%s/%s", scheme, r.Host, token, tvgId, channel.sources.MasterPlaylist())

		entry := m3uparser.M3UEntry{
			URI:   uri,
			Title: channel.sources.MediaName(),
			Tags:  make([]m3uparser.M3UTag, 0),
		}
		entry.Tags = append(entry.Tags, channel.sources.M3UTags()...)
		if !channel.sources.IsRadio() {
			entry.AddTag("KODIPROP", "inputstream=inputstream.adaptive")
			entry.AddTag("KODIPROP", "inputstream.adaptive.manifest_type=hls")
		}
		w.Write([]byte(entry.String() + "\n"))
	}
}

func (p *ChannelsHandler) manifestRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	ok := auth.VerifyToken(token)
	if !ok {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		logger.Errorf("Unauthorized access to stream stream %s: Token expired, missing, or invalid.", r.URL.Path)
		return
	}

	channelId, ok := vars["channelId"]
	if !ok {
		http.Error(w, "Invalid stream ID", http.StatusBadRequest)
		return
	}

	channel := p.GetChannel(channelId)
	if channel == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	if !channel.sources.Active() {
		http.Error(w, "Stream not active", http.StatusNotFound)
		return
	}

	channel.sources.ServeManifest(w, r, p.config.data.Timeout)
}

func (p *ChannelsHandler) mediaRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	ok := auth.VerifyToken(token)
	if !ok {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		logger.Errorf("Unauthorized access to stream stream %s: Token expired, missing, or invalid.", r.URL.Path)
		return
	}

	channelId, ok := vars["channelId"]
	if !ok {
		http.Error(w, "Invalid stream ID", http.StatusBadRequest)
		return
	}

	channel := p.GetChannel(channelId)
	if channel == nil {
		http.Error(w, "Stream not found", http.StatusNotFound)
		return
	}

	if !channel.sources.Active() {
		http.Error(w, "Stream not active", http.StatusNotFound)
		return
	}

	channel.sources.ServeMedia(w, r, p.config.data.Timeout)
}

func (p *ChannelsHandler) GetChannel(id string) *streamEntry {
	p.channelsMux.RLock()
	defer p.channelsMux.RUnlock()
	mapIndex, ok := p.channelsIdMap[id]
	if !ok {
		return nil
	}
	return p.channels[mapIndex]
}

func monitorWorker(streams <-chan *streamEntry, stop <-chan bool, wg *sync.WaitGroup) {

	defer wg.Done()
	for stream := range streams {
		select {
		case <-stop:
			return
		default:
			stream.sources.HealthCheck()
			if stream.sources.GetActiveSource() == nil {
				continue
			}
			if !stream.sources.Active() {
				logger.Warnf("Stream %s is not active", stream.sources.GetActiveSource().MediaName())
			}
		}
	}
}
