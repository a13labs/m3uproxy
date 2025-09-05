package streamserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/a13labs/m3uproxy/pkg/auth"
	authproviders "github.com/a13labs/m3uproxy/pkg/auth/providers"
	"github.com/a13labs/m3uproxy/pkg/provider"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	config      *ServerConfig
	restartChan *chan bool
	channels    *ChannelsHandler
}

func NewAPIHandler(config *ServerConfig, restartChan *chan bool, channels *ChannelsHandler) *APIHandler {

	return &APIHandler{
		config:      config,
		restartChan: restartChan,
		channels:    channels,
	}
}

func (h *APIHandler) RegisterRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/v1/authenticate", basicAuth(h.authenticateRequest))
	r.HandleFunc("/api/v1/reload", adminAccess(h.reloadRequest))
	r.HandleFunc("/api/v1/config", adminAccess(h.configAPIRequest))
	r.HandleFunc("/api/v1/playlist", adminAccess(h.playlistAPIRequest))
	r.HandleFunc("/api/v1/users", adminAccess(h.usersAPIRequest))
	r.HandleFunc("/api/v1/user/{id}", adminAccess(h.userAPIRequest))
	r.HandleFunc("/api/v1/diags/channel/{id}", adminAccess(h.diagnosticChannelRequest))
	return r
}

func (h *APIHandler) authenticateRequest(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	authParts := strings.SplitN(authHeader, " ", 2)
	token := authParts[1]

	role, err := auth.GetRoleFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := auth.GetUserFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf(`{"role": "%s", "user": "%s", "token": "%s"}`, role, user, token)
	w.Write([]byte(resp))
}

func (h *APIHandler) configAPIRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(h.config.data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	case http.MethodPut:
		newConfig := ConfigData{}
		err := json.NewDecoder(r.Body).Decode(&newConfig)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.config.Set(newConfig)
		err = h.config.Save()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) usersAPIRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		users, err := auth.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) userAPIRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		username := vars["id"]
		user, err := auth.GetUser(username)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	case http.MethodPost:
		vars := mux.Vars(r)
		username := vars["id"]
		newUser := authproviders.User{}
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = auth.AddUser(username, newUser.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newUser.Role != "" {
			err = auth.SetRole(username, newUser.Role)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodPut:
		vars := mux.Vars(r)
		username := vars["id"]
		newUser := authproviders.User{}
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newUser.Username != username {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newUser.Password != "" {
			err = auth.ChangePassword(username, newUser.Password)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if newUser.Role != "" {
			err = auth.SetRole(username, newUser.Role)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		vars := mux.Vars(r)
		username := vars["id"]
		err := auth.RemoveUser(username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) playlistAPIRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		file, err := os.Open(h.config.data.Playlist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		config := provider.PlaylistConfig{}
		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data, err := json.Marshal(config)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(data))
		return
	case http.MethodPost:
		playlist := provider.PlaylistConfig{}
		err := json.NewDecoder(r.Body).Decode(&playlist)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !playlist.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = playlist.SaveToFile(h.config.data.Playlist)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) reloadRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		w.WriteHeader(http.StatusNoContent)
		if h.restartChan != nil {
			*h.restartChan <- true
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *APIHandler) diagnosticChannelRequest(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)

	if h.channels == nil {
		return
	}

	channel := h.channels.GetChannel(id)
	if channel == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(channel.sources.Diagnostic())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(data))
}
