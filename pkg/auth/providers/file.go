package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

type FileAuthProviderConfig struct {
	FilePath string `json:"file_path"`
}

type FileAuthProvider struct {
	AuthProvider
	userStoreMux sync.Mutex
	users        Users
	config       FileAuthProviderConfig
}

func NewFileAuthProvider(config json.RawMessage) AuthProvider {

	var c FileAuthProviderConfig
	if err := json.Unmarshal([]byte(config), &c); err != nil {
		return nil
	}
	return &FileAuthProvider{config: c}
}

func (a *FileAuthProvider) AuthenticateUser(username, password string) bool {
	err := a.LoadUsers()
	if err != nil {
		return false
	}
	hashedPassword := HashPassword(password)
	for _, user := range a.users.Users {
		if user.Username == username && user.Password == hashedPassword {
			return true
		}
	}
	return false
}

func (a *FileAuthProvider) AddUser(username, password string) error {
	err := a.LoadUsers()
	if err != nil {
		return err
	}
	// check if user already exists
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for _, user := range a.users.Users {
		if user.Username == username {
			return fmt.Errorf("user already exists")
		}
	}
	a.users.Users = append(a.users.Users, User{Username: username, Password: HashPassword(password)})
	data, err := json.MarshalIndent(a.users, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.Create(a.config.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func (a *FileAuthProvider) RemoveUser(username string) error {
	err := a.LoadUsers()
	if err != nil {
		return err
	}
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for i, user := range a.users.Users {
		if user.Username == username {
			a.users.Users = append(a.users.Users[:i], a.users.Users[i+1:]...)
			data, err := json.MarshalIndent(a.users, "", "  ")
			if err != nil {
				return err
			}
			file, err := os.Create(a.config.FilePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.Write(data)
			return err
		}
	}
	return fmt.Errorf("user not found")
}

func (a *FileAuthProvider) GetUsers() ([]string, error) {
	err := a.LoadUsers()
	if err != nil {
		return nil, err
	}
	var usernames []string
	for _, user := range a.users.Users {
		usernames = append(usernames, user.Username)
	}
	return usernames, nil
}

func (a *FileAuthProvider) ChangePassword(username, password string) error {
	err := a.LoadUsers()
	if err != nil {
		return err
	}
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for i, user := range a.users.Users {
		if user.Username == username {
			a.users.Users[i].Password = HashPassword(password)
			data, err := json.MarshalIndent(a.users, "", "  ")
			if err != nil {
				return err
			}
			file, err := os.Create(a.config.FilePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.Write(data)
			return err
		}
	}
	return fmt.Errorf("user not found")
}

func (a *FileAuthProvider) DropUsers() error {
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	file, err := os.Create(a.config.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte("{\"users\":[]}"))
	return err
}

func (a *FileAuthProvider) LoadUsers() error {
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()

	if _, err := os.Stat(a.config.FilePath); os.IsNotExist(err) {
		file, err := os.Create(a.config.FilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write([]byte("{\"users\":[]}"))
		if err != nil {
			return err
		}
	}

	file, err := os.Open(a.config.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	content := string(body)

	if err := json.Unmarshal([]byte(content), &a.users); err != nil {
		return err
	}
	return nil
}

func (a *FileAuthProvider) GetRole(username string) (string, error) {
	err := a.LoadUsers()
	if err != nil {
		return "", err
	}
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for i, user := range a.users.Users {
		if user.Username == username {
			return a.users.Users[i].Role, nil
		}
	}

	return "", fmt.Errorf("user not found")
}

func (a *FileAuthProvider) GetUser(username string) (UserView, error) {
	err := a.LoadUsers()
	if err != nil {
		return UserView{}, err
	}
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for i, user := range a.users.Users {
		if user.Username == username {
			return UserView{
				Username: a.users.Users[i].Username,
				Role:     a.users.Users[i].Role,
			}, nil
		}
	}
	return UserView{}, fmt.Errorf("user not found")
}

func (a *FileAuthProvider) SetRole(username, role string) error {
	err := a.LoadUsers()
	if err != nil {
		return err
	}
	a.userStoreMux.Lock()
	defer a.userStoreMux.Unlock()
	for i, user := range a.users.Users {
		if user.Username == username {
			a.users.Users[i].Role = role
			data, err := json.MarshalIndent(a.users, "", "  ")
			if err != nil {
				return err
			}
			file, err := os.Create(a.config.FilePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.Write(data)
			return err
		}
	}
	return fmt.Errorf("user not found")
}
