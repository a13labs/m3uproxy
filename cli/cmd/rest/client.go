package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIConfig struct {
	Host     string
	Username string
	Password string
	Token    string
}

var Config APIConfig

func Authenticate() error {
	hostURL, err := url.Parse(Config.Host)
	if err != nil {
		return err
	}
	apiProtocol := hostURL.Scheme
	apiHost := hostURL.Hostname()
	apiPort := hostURL.Port()
	apiUser := Config.Username
	apiPassword := Config.Password
	apiURL := fmt.Sprintf("%s://%s:%s/api/v1/authenticate", apiProtocol, apiHost, apiPort)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(apiUser, apiPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	token, ok := result["token"]
	if !ok {
		return fmt.Errorf("token not found in response")
	}

	Config.Token = token
	return nil
}

func Call(method, endpoint string, body interface{}) (string, error) {
	var err error
	hostURL, err := url.Parse(Config.Host)
	if err != nil {
		return "", err
	}
	apiProtocol := hostURL.Scheme
	apiHost := hostURL.Hostname()
	apiPort := hostURL.Port()
	client := &http.Client{}
	apiURL := fmt.Sprintf("%s://%s:%s%s", apiProtocol, apiHost, apiPort, endpoint)

	var reqBody []byte
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return "", err
		}
	}

	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+Config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API error: %s", string(respBody))
	}

	return string(respBody), nil
}
