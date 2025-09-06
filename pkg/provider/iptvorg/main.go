package iptvorg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	types "github.com/a13labs/m3uproxy/pkg/provider/types"
)

const (
	IPTV_API_URL = "https://iptv-org.github.io/api"
)

type IPTVOrgProvider struct {
	types.M3UProvider
	playlist m3uparser.M3UPlaylist
}

const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

type IPTVOrgConfig struct {
	Categories []string `json:"categories"`
	Countries  []string `json:"countries"`
	UserAgent  string   `json:"user_agent,omitempty"`
}

type IPTVOrgChannel struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Country    string   `json:"country"`
	Categories []string `json:"categories"`
	Website    string   `json:"website,omitempty"`
	Logo       string   `json:"logo,omitempty"`
}

type IPTVOrgStream struct {
	Channel      string `json:"channel"`
	URL          string `json:"url"`
	Timeshift    string `json:"timeshift,omitempty"`
	HTTPReferrer string `json:"http_referrer,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
}

type cachedEntry struct {
	m3uEntry    *m3uparser.M3UEntry
	iptvChannel *IPTVOrgChannel
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getChannels(config IPTVOrgConfig) (map[string]cachedEntry, error) {

	resp, err := http.Get(IPTV_API_URL + "/channels.json")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch channels: %s", resp.Status)
	}

	defer resp.Body.Close()

	remoteChannels := []IPTVOrgChannel{}
	err = json.NewDecoder(resp.Body).Decode(&remoteChannels)

	if err != nil {
		return nil, err
	}

	var channels = make(map[string]cachedEntry)

	for i, channel := range remoteChannels {
		inCategories := len(config.Categories) == 0
		inCountries := len(config.Countries) == 0
		for _, category := range config.Categories {
			inCategories = inCategories || contains(channel.Categories, category)
		}
		for _, country := range config.Countries {
			inCountries = inCountries || channel.Country == country
		}
		if inCategories && inCountries {
			channels[channel.ID] = cachedEntry{
				iptvChannel: &remoteChannels[i],
				m3uEntry:    nil,
			}
		}
	}

	return channels, nil
}

func getStreams(channels map[string]cachedEntry, config IPTVOrgConfig) ([]m3uparser.M3UEntry, error) {

	resp, err := http.Get(IPTV_API_URL + "/streams.json")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	streams := []IPTVOrgStream{}
	err = json.NewDecoder(resp.Body).Decode(&streams)

	if err != nil {
		return nil, err
	}

	entries := make(m3uparser.M3UEntries, 0)
	for id, cache := range channels {
		for _, stream := range streams {

			if stream.Channel == "" {
				continue
			}

			if stream.Channel == id {
				if cache.m3uEntry != nil {
					continue
				}

				if stream.UserAgent == "" {
					stream.UserAgent = config.UserAgent
				}

				if stream.HTTPReferrer == "" {
					referrer, err := url.Parse(cache.iptvChannel.Website)
					if err == nil {
						if referrer.Scheme == "" {
							referrer.Scheme = "http"
						}
						if referrer.Host != "" {
							stream.HTTPReferrer = referrer.Scheme + "://" + referrer.Host
						}
					}
				}

				headers := make(map[string]string)
				if stream.HTTPReferrer != "" {
					headers["http-referer"] = stream.HTTPReferrer
				}

				headers["http-user-agent"] = stream.UserAgent

				extinftags := make(m3uparser.M3UExtinfTags, 0)
				extinftags = append(extinftags, m3uparser.M3UTvgTag{
					Tag:   "tvg-id",
					Value: cache.iptvChannel.ID,
				})
				extinftags = append(extinftags, m3uparser.M3UTvgTag{
					Tag:   "tvg-name",
					Value: cache.iptvChannel.Name,
				})
				extinftags = append(extinftags, m3uparser.M3UTvgTag{
					Tag:   "tvg-logo",
					Value: cache.iptvChannel.Logo,
				})
				extinftags = append(extinftags, m3uparser.M3UTvgTag{
					Tag:   "tvg-country",
					Value: cache.iptvChannel.Country,
				})
				extinftags = append(extinftags, m3uparser.M3UTvgTag{
					Tag:   "tvg-group",
					Value: "TV",
				})
				if len(cache.iptvChannel.Categories) != 0 {
					extinftags = append(extinftags, m3uparser.M3UTvgTag{
						Tag:   "tvg-type",
						Value: cache.iptvChannel.Categories[0],
					})
				}

				tags := make(m3uparser.M3UTags, 0)
				tags = append(tags, m3uparser.M3UTag{
					Tag:   "EXTINF",
					Value: fmt.Sprintf("-1 %s, %s", extinftags.String(), cache.iptvChannel.Name),
				})

				for k, v := range headers {
					tags = append(tags, m3uparser.M3UTag{
						Tag:   "EXTVLCOPT",
						Value: fmt.Sprintf("%s=%s", k, v),
					})
				}

				cache.m3uEntry = &m3uparser.M3UEntry{
					Title:      cache.iptvChannel.Name,
					URI:        stream.URL,
					Tags:       tags,
					ExtInfTags: extinftags,
				}

				entries = append(entries, *cache.m3uEntry)
			}
		}
	}

	return entries, nil
}

func (p *IPTVOrgProvider) GetPlaylist() *m3uparser.M3UPlaylist {
	return &p.playlist
}

func NewIPTVOrgProvider(config json.RawMessage) *IPTVOrgProvider {

	cfg := IPTVOrgConfig{}
	err := json.Unmarshal([]byte(config), &cfg)
	if err != nil {
		return nil
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = DefaultUserAgent
	}

	channels, err := getChannels(cfg)
	if err != nil {
		return nil
	}

	streams, err := getStreams(channels, cfg)
	if err != nil {
		return nil
	}

	return &IPTVOrgProvider{
		playlist: m3uparser.M3UPlaylist{
			Entries: streams,
		},
	}
}
