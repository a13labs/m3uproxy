package auth

import (
	"encoding/json"
	"errors"

	"github.com/a13labs/m3uproxy/pkg/auth/providers"
)

type AuthConfig struct {
	Provider       string          `json:"provider"`
	SecretKey      string          `json:"secret_key"`
	ExpirationTime int             `json:"expiration_time,omitempty"`
	Settings       json.RawMessage `json:"settings"`
}

var authConfig AuthConfig

func InitializeAuth(data json.RawMessage) error {

	err := json.Unmarshal(data, &authConfig)
	if err != nil {
		return err
	}

	if authConfig.Provider == "" {
		return errors.New("auth provider is required")
	}

	if len(authConfig.SecretKey) == 0 {
		return errors.New("secret key is required")
	}

	if authConfig.ExpirationTime == 0 {
		authConfig.ExpirationTime = 24
	}

	return providers.InitializeAuthProvider(authConfig.Provider, authConfig.Settings)
}

func GetRole(username string) (string, error) {
	return providers.GetRole(username)
}

func CheckCredentials(username, password string) bool {
	return providers.AuthenticateUser(username, password)
}

func AddUser(username, password string) error {
	return providers.AddUser(username, password)
}

func RemoveUser(username string) error {
	return providers.RemoveUser(username)
}

func GetUsers() ([]string, error) {
	return providers.GetUsers()
}

func ChangePassword(username, password string) error {
	return providers.ChangePassword(username, password)
}

func DropUsers() error {
	return providers.DropUsers()
}

func GetUser(username string) (providers.UserView, error) {
	return providers.GetUser(username)
}

func SetRole(username, role string) error {
	return providers.SetRole(username, role)
}
