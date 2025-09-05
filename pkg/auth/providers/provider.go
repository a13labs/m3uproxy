package providers

import (
	"encoding/json"
	"fmt"
)

type AuthProvider interface {
	AuthenticateUser(username, password string) bool
	AddUser(username, password string) error
	RemoveUser(username string) error
	GetRole(username string) (string, error)
	SetRole(username, role string) error
	GetUsers() ([]string, error)
	ChangePassword(username, password string) error
	DropUsers() error
	LoadUsers() error
	GetUser(username string) (UserView, error)
}

type AuthProviderFactory func(config json.RawMessage) AuthProvider

var authProvider AuthProvider

func SetAuthProviderFactory(factory AuthProviderFactory, config json.RawMessage) {
	authProvider = factory(config)
}

func GetAuthProvider() AuthProvider {
	return authProvider
}

func AuthenticateUser(username, password string) bool {
	return authProvider.AuthenticateUser(username, password)
}

func AddUser(username, password string) error {
	return authProvider.AddUser(username, password)
}

func RemoveUser(username string) error {
	return authProvider.RemoveUser(username)
}

func GetUsers() ([]string, error) {
	return authProvider.GetUsers()
}

func ChangePassword(username, password string) error {
	return authProvider.ChangePassword(username, password)
}

func DropUsers() error {
	return authProvider.DropUsers()
}

func GetRole(username string) (string, error) {
	return authProvider.GetRole(username)
}

func InitializeAuthProvider(provider string, config json.RawMessage) error {
	switch provider {
	case "file":
		SetAuthProviderFactory(NewFileAuthProvider, config)
	case "memory":
		SetAuthProviderFactory(NewMemoryAuthProvider, config)
	case "null":
		SetAuthProviderFactory(NewNullAuthProvider, config)
	default:
		return fmt.Errorf("unsupported auth provider")
	}
	return nil
}

func LoadUsers() error {
	return authProvider.LoadUsers()
}

func GetUser(username string) (UserView, error) {
	return authProvider.GetUser(username)
}

func SetRole(username, role string) error {
	return authProvider.SetRole(username, role)
}
