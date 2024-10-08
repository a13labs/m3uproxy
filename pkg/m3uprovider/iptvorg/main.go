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
package iptvorg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	types "github.com/a13labs/m3uproxy/pkg/m3uprovider/types"
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
		log.Printf("Error getting channels: %s", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error getting channels: %s", resp.Status)
		return nil, err
	}

	defer resp.Body.Close()

	remoteChannels := []IPTVOrgChannel{}
	err = json.NewDecoder(resp.Body).Decode(&remoteChannels)

	if err != nil {
		log.Printf("Error decoding channels: %s", err)
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
			log.Printf("Adding channel %s, id: %s", channel.Name, channel.ID)
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
		log.Printf("Error getting streams: %s", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error getting streams: %s", resp.Status)
		return nil, err
	}

	defer resp.Body.Close()
	streams := []IPTVOrgStream{}
	err = json.NewDecoder(resp.Body).Decode(&streams)

	if err != nil {
		log.Printf("Error decoding streams: %s", err)
		return nil, err
	}

	entries := make(m3uparser.M3UEntries, 0)
	for id, cache := range channels {
		log.Printf("Channel %s, id: %s", cache.iptvChannel.Name, id)
		for _, stream := range streams {

			if stream.Channel == "" {
				continue
			}

			if stream.Channel == id {
				if cache.m3uEntry != nil {
					log.Printf("Stream for channel %s already added, url: %s", cache.iptvChannel.Name, stream.URL)
					continue
				}

				log.Printf("Checking stream for channel %s, url: %s", cache.iptvChannel.Name, stream.URL)

				if stream.UserAgent == "" {
					stream.UserAgent = config.UserAgent
				}

				if stream.HTTPReferrer == "" {
					referrer, err := url.Parse(cache.iptvChannel.Website)
					if err != nil {
						log.Printf("Error parsing URL: %s", err)
					} else {
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

				tvgtags := make(m3uparser.M3UTvgTags, 0)
				tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
					Tag:   "tvg-id",
					Value: cache.iptvChannel.ID,
				})
				tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
					Tag:   "tvg-name",
					Value: cache.iptvChannel.Name,
				})
				tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
					Tag:   "tvg-logo",
					Value: cache.iptvChannel.Logo,
				})
				tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
					Tag:   "tvg-country",
					Value: cache.iptvChannel.Country,
				})
				tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
					Tag:   "tvg-group",
					Value: "TV",
				})
				if cache.iptvChannel.Categories != nil {
					tvgtags = append(tvgtags, m3uparser.M3UTvgTag{
						Tag:   "tvg-type",
						Value: cache.iptvChannel.Categories[0],
					})
				}

				tags := make(m3uparser.M3UTags, 0)
				tags = append(tags, m3uparser.M3UTag{
					Tag:   "EXTINF",
					Value: fmt.Sprintf("-1 %s, %s", tvgtags.String(), cache.iptvChannel.Name),
				})

				for k, v := range headers {
					tags = append(tags, m3uparser.M3UTag{
						Tag:   "EXTVLCOPT",
						Value: fmt.Sprintf("%s=%s", k, v),
					})
				}

				cache.m3uEntry = &m3uparser.M3UEntry{
					Title:   cache.iptvChannel.Name,
					URI:     stream.URL,
					Tags:    tags,
					TVGTags: tvgtags,
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
		log.Println("Error parsing config")
		return nil
	}

	log.Println("Config:")
	log.Printf("Categories: %v", cfg.Categories)
	log.Printf("Countries: %v", cfg.Countries)
	if cfg.UserAgent == "" {
		cfg.UserAgent = DefaultUserAgent
	}

	log.Println("Getting channels from iptv.org")
	channels, err := getChannels(cfg)
	if err != nil {
		log.Println("Error getting channels")
		return nil
	}

	log.Println("Getting streams from iptv.org")
	streams, err := getStreams(channels, cfg)
	if err != nil {
		log.Println("Error getting streams")
		return nil
	}

	return &IPTVOrgProvider{
		playlist: m3uparser.M3UPlaylist{
			Entries: streams,
		},
	}
}
