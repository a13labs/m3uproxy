/*
Copyright © 2024 Alexandre Pires

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package userstore

import (
	"encoding/json"

	"github.com/a13labs/m3uproxy/pkg/userstore/auth"
)

func InitializeAuthProvider(provider string, config json.RawMessage) error {
	return auth.InitializeAuthProvider(provider, config)
}

func AuthenticateUser(username, password string) bool {
	return auth.AuthenticateUser(username, password)
}

func AddUser(username, password string) error {
	return auth.AddUser(username, password)
}

func RemoveUser(username string) error {
	return auth.RemoveUser(username)
}

func GetUsers() ([]string, error) {
	return auth.GetUsers()
}

func ChangePassword(username, password string) error {
	return auth.ChangePassword(username, password)
}

func DropUsers() error {
	return auth.DropUsers()
}
