package streamserver

import (
	"net/http"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/gorilla/mux"
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
