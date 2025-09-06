package streamserver

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/gorilla/mux"

	"github.com/a13labs/m3uproxy/pkg/xmltv"
)

type EPGHandler struct {
	config *ServerConfig
}

func NewEPGHandler(config *ServerConfig) *EPGHandler {
	return &EPGHandler{
		config: config,
	}
}

func (e *EPGHandler) RegisterRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/epg.xml", basicAuth(e.epgRequest))
	r.HandleFunc("/epg/{channel_id}", basicAuth(e.channelEpgRequest))
	return r
}

func (e *EPGHandler) epgRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	content, err := loadContent(e.config.data.Epg)
	if err != nil {
		http.Error(w, "EPG file not found", http.StatusNotFound)
		logger.Errorf("EPG file not found at %s", e.config.data.Epg)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

func (e *EPGHandler) channelEpgRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	channelID := vars["channel_id"]
	timestamp := r.URL.Query().Get("timestamp")
	if timestamp != "" {
		_, err := parseTimestamp(timestamp)
		if err != nil {
			http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
			return
		}
	}

	content, err := loadContent(e.config.data.Epg)
	if err != nil {
		logger.Errorf("EPG file not found at %s", e.config.data.Epg)
		http.Error(w, "EPG file not found", http.StatusNotFound)
		return
	}

	// Parse XMLTV data
	var epg xmltv.EPG
	if err := xml.Unmarshal([]byte(content), &epg); err != nil {
		logger.Errorf("Error parsing XMLTV: %v", err)
		http.Error(w, "Error parsing EPG", http.StatusInternalServerError)
		return
	}

	// Access the data
	logger.Infof("Found %d channels and %d programmes\n", len(epg.Channels), len(epg.Programmes))
	filteredEPG := xmltv.EPG{
		Channels:   []xmltv.Channel{},
		Programmes: []xmltv.Programme{},
	}

	for _, channel := range epg.Channels {
		if channel.ID == channelID {
			filteredEPG.Channels = append(filteredEPG.Channels, channel)
			break
		}
	}

	for _, programme := range epg.Programmes {
		if programme.Channel == channelID {
			if timestamp != "" {
				startTime, err := parseTimestamp(programme.Start.Format("20060102150405 -0700"))
				if err != nil {
					http.Error(w, "Error parsing programme start time", http.StatusInternalServerError)
					logger.Errorf("Error parsing programme start time: %v", err)
					return
				}
				endTime, err := parseTimestamp(programme.Stop.Format("20060102150405 -0700"))
				if err != nil {
					http.Error(w, "Error parsing programme stop time", http.StatusInternalServerError)
					logger.Errorf("Error parsing programme stop time: %v", err)
					return
				}
				reqTime, err := parseTimestamp(timestamp)
				if err != nil {
					http.Error(w, "Error parsing request timestamp", http.StatusInternalServerError)
					logger.Errorf("Error parsing request timestamp: %v", err)
					return
				}
				if reqTime.Before(*startTime) || reqTime.After(*endTime) {
					continue
				}
			}
			filteredEPG.Programmes = append(filteredEPG.Programmes, programme)
		}
	}

	filteredXML, err := xml.MarshalIndent(filteredEPG, "", "  ")
	if err != nil {
		http.Error(w, "Error generating filtered EPG", http.StatusInternalServerError)
		logger.Errorf("Error generating filtered EPG: %v", err)
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
