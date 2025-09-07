package streamserver

import (
	"encoding/json"
	"os"
)

type GeoIPConfig struct {
	Database         string   `json:"database"`
	Whitelist        []string `json:"whitelist,omitempty"`
	InternalNetworks []string `json:"internal_networks,omitempty"`
}

type SecurityConfig struct {
	GeoIP              GeoIPConfig `json:"geoip,omitempty"`
	AllowedCORSDomains []string    `json:"allowed_cors_domains,omitempty"`
}

type ConfigData struct {
	Port        int             `json:"port"`
	Playlist    string          `json:"playlist"`
	Epg         []string        `json:"epg"`
	Timeout     int             `json:"default_timeout,omitempty"`
	NumWorkers  int             `json:"num_workers,omitempty"`
	ScanTime    int             `json:"scan_time,omitempty"`
	EPGScanTime int             `json:"epg_scan_time,omitempty"`
	Security    SecurityConfig  `json:"security,omitempty"`
	Auth        json.RawMessage `json:"auth"`
	LogFile     string          `json:"log_file,omitempty"`
	LogLevel    string          `json:"log_level,omitempty"`
	CacheDir    string          `json:"cache_dir,omitempty"`
}

type ServerConfig struct {
	path string
	data ConfigData
}

func NewServerConfig(path string) *ServerConfig {
	c := &ServerConfig{
		path: path,
	}

	if err := c.Load(path); err != nil {
		if os.IsNotExist(err) {
			c.data = ConfigData{
				Port:        8080,
				Playlist:    "playlist.m3u",
				Epg:         []string{"epg.xml"},
				Timeout:     5,
				NumWorkers:  4,
				ScanTime:    60,
				EPGScanTime: 86400,
				Security: SecurityConfig{
					GeoIP: GeoIPConfig{
						Database:         "GeoLite2-Country.mmdb",
						Whitelist:        []string{},
						InternalNetworks: []string{},
					},
					AllowedCORSDomains: []string{},
				},
				Auth:     json.RawMessage("{}"),
				LogFile:  "server.log",
				CacheDir: "cache",
			}
			if err := c.Save(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	// Ensure CacheDir exists
	if c.data.CacheDir == "" {
		c.data.CacheDir = "cache"
	}
	if _, err := os.Stat(c.data.CacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(c.data.CacheDir, 0755); err != nil {
			panic(err)
		}
	}
	return c
}

func (c *ConfigData) Merge(other ConfigData) {
	if other.Timeout != 0 {
		c.Timeout = other.Timeout
	}
	if other.NumWorkers != 0 {
		c.NumWorkers = other.NumWorkers
	}
	if other.ScanTime != 0 {
		c.ScanTime = other.ScanTime
	}
	if len(other.Security.GeoIP.Whitelist) > 0 {
		c.Security.GeoIP.Whitelist = other.Security.GeoIP.Whitelist
	}
	if len(other.Security.GeoIP.InternalNetworks) > 0 {
		c.Security.GeoIP.InternalNetworks = other.Security.GeoIP.InternalNetworks
	}
}

func (c *ServerConfig) Load(path string) error {

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&c.data)
	if err != nil {
		return err
	}

	c.path = path

	return nil
}

func (c *ServerConfig) Get() ConfigData {
	return c.data
}

func (c *ServerConfig) Set(data ConfigData) {
	c.data = data
}

func (c *ServerConfig) GetPath() string {
	return c.path
}

func (c *ServerConfig) SetPath(path string) {
	c.path = path
}

func (c *ServerConfig) GetPlaylist() string {
	return c.data.Playlist
}

func (c *ServerConfig) SetPlaylist(playlist string) {
	c.data.Playlist = playlist
}

func (c *ServerConfig) GetEpg() []string {
	return c.data.Epg
}

func (c *ServerConfig) SetEpg(epg []string) {
	c.data.Epg = epg
}

func (c *ServerConfig) GetTimeout() int {
	return c.data.Timeout
}

func (c *ServerConfig) SetTimeout(timeout int) {
	c.data.Timeout = timeout
}

func (c *ServerConfig) GetNumWorkers() int {
	return c.data.NumWorkers
}

func (c *ServerConfig) SetNumWorkers(numWorkers int) {
	c.data.NumWorkers = numWorkers
}

func (c *ServerConfig) GetScanTime() int {
	return c.data.ScanTime
}

func (c *ServerConfig) SetScanTime(scanTime int) {
	c.data.ScanTime = scanTime
}

func (c *ServerConfig) GetSecurity() SecurityConfig {
	return c.data.Security
}

func (c *ServerConfig) SetSecurity(security SecurityConfig) {
	c.data.Security = security
}

func (c *ServerConfig) GetAuth() json.RawMessage {
	return c.data.Auth
}

func (c *ServerConfig) SetAuth(auth json.RawMessage) {
	c.data.Auth = auth
}

func (c *ServerConfig) Save() error {

	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.data)
}
