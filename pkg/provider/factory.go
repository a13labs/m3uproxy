package provider

import (
	"errors"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	"github.com/a13labs/m3uproxy/pkg/provider/file"
	"github.com/a13labs/m3uproxy/pkg/provider/iptvorg"
	types "github.com/a13labs/m3uproxy/pkg/provider/types"
)

func NewProvider(config ProviderConfig) types.M3UProvider {

	switch config.Provider {
	case "iptv.org":
		return iptvorg.NewIPTVOrgProvider(config.Config)
	case "file":
		return file.NewM3UFileProvider(config.Config)
	default:
		return nil
	}
}

func Load(config *PlaylistConfig) (*m3uparser.M3UPlaylist, error) {

	providersPriority := make([]string, 0)
	if config.ProvidersPriority != nil {
		if len(config.ProvidersPriority) != len(config.Providers) {
			return nil, errors.New("providers_priority and providers must have the same length")
		}
		providersPriority = append(providersPriority, config.ProvidersPriority...)
	} else {
		for providerName := range config.Providers {
			providersPriority = append(providersPriority, providerName)
		}
	}

	masterPlaylist := m3uparser.M3UPlaylist{
		Version: 3,
		Entries: make(m3uparser.M3UEntries, 0),
		Tags:    make(m3uparser.M3UTags, 0),
	}

	for _, providerName := range providersPriority {

		provider := NewProvider(config.Providers[providerName])
		if provider == nil {
			return nil, errors.New("provider not available '" + providerName + "'")
		}

		playlist := provider.GetPlaylist()
		ignoreTags := config.Providers[providerName].IgnoreTags
		for _, entry := range playlist.Entries {

			skip := false
			for _, tag := range entry.ExtInfTags {
				v, ok := ignoreTags[tag.Tag]
				skip = skip || (ok && v == tag.Value)
			}

			if skip {
				logger.Infof("Channel '%s' is ignored, skipping.", entry.Title)
				continue
			}

			tvgId := entry.ExtInfTags.GetValue("tvg-id")
			if tvgId == "" {
				tvgId = entry.Title
			}

			override, ok := config.Overrides[tvgId]
			if ok && override.Disabled {
				logger.Infof("Channel '%s' is disabled, skipping.", entry.Title)
				continue
			}
			if ok && override.ChannelName != "" {
				entry.Title = override.ChannelName
			}
			if ok && override.URL != "" {
				entry.URI = override.URL
			}
			if ok && len(override.Headers) > 0 {
				for k, v := range override.Headers {
					entry.Tags = append(entry.Tags, m3uparser.M3UTag{
						Tag:   "M3UPROXYHEADER",
						Value: k + "=" + v,
					})
				}
			}
			if ok && override.HttpProxy != "" {
				entry.Tags = append(entry.Tags, m3uparser.M3UTag{
					Tag:   "M3UPROXYTRANSPORT",
					Value: "proxy=" + override.HttpProxy,
				})
			}
			if ok && override.ForceKodiHeaders {
				entry.Tags = append(entry.Tags, m3uparser.M3UTag{
					Tag:   "M3UPROXYOPT",
					Value: "forcekodiheaders",
				})
			}
			if ok && override.DisableRemap {
				entry.Tags = append(entry.Tags, m3uparser.M3UTag{
					Tag:   "M3UPROXYOPT",
					Value: "disableremap",
				})
			}
			masterPlaylist.Entries = append(masterPlaylist.Entries, entry)
		}
	}

	if len(config.ChannelOrder) > 0 {
		logger.Info("Ordering playlist by provided channel order.")

		for needle, channel := range config.ChannelOrder {
			for pos := needle; pos < len(masterPlaylist.Entries); pos++ {
				if masterPlaylist.Entries[pos].ExtInfTags.GetValue("tvg-id") == channel {
					if needle == pos {
						break
					}
					masterPlaylist.Entries[needle], masterPlaylist.Entries[pos] = masterPlaylist.Entries[pos], masterPlaylist.Entries[needle]
					break
				}
			}
		}
	}

	return &masterPlaylist, nil
}
