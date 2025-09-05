package streamserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/a13labs/m3uproxy/pkg/auth"
	"github.com/a13labs/m3uproxy/pkg/logger"
)

type streamLicense struct {
	keyType string
	keyId   string
	key     string
}

type streamLicenseType struct {
	keys map[string]streamLicense
}

type streamLicenseManager struct {
	licenses map[string]streamLicenseType
}

type emeLicenseKey struct {
	Type string `json:"kty"`
	Key  string `json:"k"`
	Id   string `json:"kid"`
}

type emeLicenseResponse struct {
	Keys []emeLicenseKey `json:"keys"`
}

func newStreamLicenseManager() *streamLicenseManager {
	return &streamLicenseManager{
		licenses: make(map[string]streamLicenseType),
	}
}

func (m *streamLicenseManager) addLicense(keyType, keyId, key string) {

	if _, ok := m.licenses[keyType]; !ok {
		m.licenses[keyType] = streamLicenseType{
			keys: make(map[string]streamLicense),
		}
	}

	m.licenses[keyType].keys[keyId] = streamLicense{
		keyType: keyType,
		keyId:   keyId,
		key:     key,
	}
}

func licenseKeysRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	authParts := strings.SplitN(authHeader, " ", 2)
	token := authParts[1]

	ok := auth.VerifyToken(token)
	if !ok {
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		logger.Errorf("Unauthorized access to stream stream %s: Token expired, missing, or invalid.", r.URL.Path)
		return
	}

	// Check if the license manager is initialized
	if licenseManger == nil {
		http.Error(w, "No licenses found", http.StatusNotFound)
		return
	}

	var licenses emeLicenseResponse
	for _, license := range licenseManger.licenses["clearkey"].keys {

		licenses.Keys = append(licenses.Keys, emeLicenseKey{
			Type: "oct",
			Key:  license.key,
			Id:   license.keyId,
		})
	}

	// Return the license
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(licenses)
}
