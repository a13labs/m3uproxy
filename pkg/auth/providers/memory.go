package providers

import (
	"encoding/json"
	"fmt"
)

type MemoryAuthProvide struct {
	AuthProvider
	users map[string]*User
}

func NewMemoryAuthProvider(config json.RawMessage) AuthProvider {
	return &MemoryAuthProvide{
		users: make(map[string]*User),
	}
}

func (a *MemoryAuthProvide) AuthenticateUser(username, password string) bool {
	if user, ok := a.users[username]; ok {
		return user.Password == HashPassword(password)
	}
	return false
}

func (a *MemoryAuthProvide) AddUser(username, password string) error {

	if _, ok := a.users[username]; ok {
		return fmt.Errorf("user already exists")
	}

	a.users[username].Password = HashPassword(password)
	return nil
}

func (a *MemoryAuthProvide) RemoveUser(username string) error {

	if _, ok := a.users[username]; !ok {
		return fmt.Errorf("user does not exist")
	}

	delete(a.users, username)
	return nil
}

func (a *MemoryAuthProvide) GetUsers() ([]string, error) {
	users := make([]string, 0, len(a.users))
	for user := range a.users {
		users = append(users, user)
	}
	return users, nil
}

func (a *MemoryAuthProvide) ChangePassword(username, password string) error {
	if _, ok := a.users[username]; !ok {
		return fmt.Errorf("user does not exist")
	}

	a.users[username].Password = HashPassword(password)
	return nil
}

func (a *MemoryAuthProvide) DropUsers() error {
	a.users = make(map[string]*User)
	return nil
}

func (a *MemoryAuthProvide) LoadUsers() error {
	return fmt.Errorf("not implemented")
}

func (a *MemoryAuthProvide) GetRole(username string) (string, error) {
	if user, ok := a.users[username]; ok {
		return user.Role, nil
	}
	return "", fmt.Errorf("user does not exist")
}

func (a *MemoryAuthProvide) GetUser(username string) (UserView, error) {
	if user, ok := a.users[username]; ok {
		return UserView{
			Username: user.Username,
			Role:     user.Role,
		}, nil
	}
	return UserView{}, fmt.Errorf("user does not exist")
}

func (a *MemoryAuthProvide) SetRole(username, role string) error {
	if user, ok := a.users[username]; ok {
		user.Role = role
		return nil
	}
	return fmt.Errorf("user does not exist")
}
