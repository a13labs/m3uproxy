package provider

import (
	"errors"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	"github.com/a13labs/m3uproxy/pkg/provider/file"
	"github.com/a13labs/m3uproxy/pkg/provider/iptvorg"
	types "github.com/a13labs/m3uproxy/pkg/provider/types"
)

// NewProvider returns a concrete M3U provider implementation based on the
// supplied configuration. It returns nil when the requested provider is not
// available; callers should handle that case and surface a helpful error.
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

// Load merges playlists from the configured providers into a single
// master playlist. Configuration controls ordering, per-provider tag
// filtering and per-channel overrides.
//
// Behaviour summary:
//   - If ProvidersPriority is provided, provider names must match the
//     Providers map length and that order is used; otherwise providers are
//     iterated in map order.
//   - Per-provider `IgnoreTags` can exclude entries that match tag/value
//     combinations.
//   - Per-channel `Overrides` allow renaming, URL substitution and adding
//     headers or transport hints.
func Load(config *PlaylistConfig) (*m3uparser.M3UPlaylist, error) {

	// Build provider iteration order. When ProvidersPriority is provided we
	// validate it matches the number of declared providers to avoid
	// accidental misconfiguration.
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

			// Determine if this entry should be skipped based on configured
			// ignore tag/value pairs for the current provider.
			skip := false
			for _, tag := range entry.ExtInfTags {
				v, ok := ignoreTags[tag.Tag]
				skip = skip || (ok && v == tag.Value)
			}

			if skip {
				logger.Debugf("Ignoring channel '%s' from provider '%s' (matched ignore tags)", entry.Title, providerName)
				continue
			}

			// Use tvg-id when present, otherwise fall back to the human
			// readable title. tvg-id is used as the canonical channel key.
			tvgId := entry.ExtInfTags.GetValue("tvg-id")
			if tvgId == "" {
				tvgId = entry.Title
			}

			// Apply per-channel overrides (disabled flag, rename, URL, etc.)
			override, ok := config.Overrides[tvgId]
			if ok && override.Disabled {
				logger.Debugf("Channel '%s' (tvg-id=%s) is disabled by override, skipping.", entry.Title, tvgId)
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

	// If a specific channel order is provided, reorder the resulting entries
	// so the entries listed in ChannelOrder appear first in that sequence.
	if len(config.ChannelOrder) > 0 {
		logger.Info("Reordering playlist according to configured ChannelOrder")

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
