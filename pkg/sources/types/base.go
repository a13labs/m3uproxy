package types

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/a13labs/m3uproxy/pkg/m3uparser"
	"github.com/elnormous/contenttype"
)

func (s *BaseStreamSource) HealthCheck() error {
	logger.Debugf("Starting health check for stream: %s (%s)", s.MediaName(), s.m3u.URI)

	uri, _, err := s.conn.Check("GET", s.m3u.URI)
	if err != nil {
		logger.Errorf("Health check failed for stream %s: %v", s.MediaName(), err)
		return err
	}

	s.mux.Lock()
	if uri != s.m3u.URI {
		logger.Infof("Stream %s URI updated from %s to %s", s.MediaName(), s.m3u.URI, uri)
		s.m3u.URI = uri
	}
	s.mux.Unlock()

	_, err = s.verify(s.m3u.URI)

	s.mux.Lock()
	oldActive := s.active
	s.active = err == nil
	if oldActive != s.active {
		if s.active {
			logger.Infof("Stream %s is now ACTIVE", s.MediaName())
		} else {
			logger.Warnf("Stream %s is now INACTIVE: %v", s.MediaName(), err)
		}
	} else {
		logger.Debugf("Stream %s health check completed, status unchanged: active=%t", s.MediaName(), s.active)
	}
	s.mux.Unlock()

	if err != nil {
		logger.Debugf("Health check verification failed for stream %s: %v", s.MediaName(), err)
	} else {
		logger.Debugf("Health check passed for stream %s", s.MediaName())
	}

	return err
}

func (s *BaseStreamSource) Diagnostic() StreamSourceDiag {
	logger.Debugf("Starting diagnostic for stream: %s (%s)", s.MediaName(), s.m3u.URI)

	uri, _, err := s.conn.Check("GET", s.m3u.URI)
	diag := StreamSourceDiag{
		Entry:       s.m3u,
		Headers:     s.headers,
		HttpProxy:   s.httpProxy,
		Active:      false,
		Diagnostics: make([]HttpDiags, 0),
	}
	if err != nil {
		logger.Errorf("Diagnostic connection check failed for stream %s: %v", s.MediaName(), err)
		diag.Diagnostics = append(diag.Diagnostics, HttpDiags{
			Url:    uri,
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return diag
	}

	s.mux.Lock()
	if uri != s.m3u.URI {
		logger.Infof("Diagnostic found URI redirect for stream %s: %s -> %s", s.MediaName(), s.m3u.URI, uri)
		s.m3u.URI = uri
	}
	s.mux.Unlock()

	logger.Debugf("Running detailed verification diagnostics for stream %s", s.MediaName())
	s.verifyWithDiags(s.m3u.URI, &diag)

	logger.Infof("Diagnostic completed for stream %s: active=%t, diagnostics_count=%d",
		s.MediaName(), diag.Active, len(diag.Diagnostics))

	return diag
}

func (s *BaseStreamSource) Active() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	active := s.active
	logger.Tracef("Stream %s active status queried: %t", s.MediaName(), active)
	return active
}

func (s *BaseStreamSource) MediaType() contenttype.MediaType {
	s.mux.RLock()
	defer s.mux.RUnlock()
	mediaType := s.mediaType
	logger.Tracef("Stream %s media type queried: %s", s.MediaName(), mediaType.String())
	return mediaType
}

func (s *BaseStreamSource) MediaName() string {
	return s.m3u.Title
}

func (s *BaseStreamSource) M3UTags() m3uparser.M3UTags {
	return s.m3u.Tags
}

func (s *BaseStreamSource) IsRadio() bool {
	return s.radio
}

func (s *BaseStreamSource) Url() string {
	url := s.m3u.URI
	logger.Tracef("Stream %s URL queried: %s", s.MediaName(), url)
	return url
}

func (s *BaseStreamSource) verify(mediaURI string) (contenttype.MediaType, error) {
	logger.Debugf("Verifying media URI for stream %s: %s", s.MediaName(), mediaURI)

	s.mux.RLock()
	body, _, ct, err := s.conn.Get("GET", mediaURI)
	s.mux.RUnlock()
	if err != nil {
		logger.Errorf("Failed to fetch media from %s for stream %s: %v", mediaURI, s.MediaName(), err)
		return contenttype.MediaType{}, err
	}

	logger.Debugf("Retrieved media for stream %s: content_type=%s, body_size=%d bytes",
		s.MediaName(), ct.String(), len(body))

	if !ct.MatchesAny(supportedMediaTypes...) {
		logger.Warnf("Unsupported content type %s for stream %s at %s", ct.String(), s.MediaName(), mediaURI)
		return contenttype.MediaType{}, errors.New("invalid content type")
	}

	if ct.Subtype == "vnd.apple.mpegurl" || ct.Subtype == "x-mpegurl" {
		logger.Debugf("Processing M3U playlist for stream %s", s.MediaName())
		m3uPlaylist, err := m3uparser.DecodeFromReader(bytes.NewReader(body))
		if err != nil {
			logger.Errorf("Failed to parse M3U playlist for stream %s: %v, but content type is valid.", s.MediaName(), err)
			return contenttype.MediaType{}, nil
		}

		if len(m3uPlaylist.Entries) == 0 {
			logger.Warnf("Empty M3U playlist for stream %s at %s", s.MediaName(), mediaURI)
			return contenttype.MediaType{}, errors.New("empty playlist")
		}

		logger.Debugf("M3U playlist for stream %s contains %d entries, following first entry",
			s.MediaName(), len(m3uPlaylist.Entries))

		uri, _ := url.Parse(m3uPlaylist.Entries[0].URI)

		if uri.Scheme == "" {
			logger.Debugf("Resolving relative URL for stream %s: %s", s.MediaName(), m3uPlaylist.Entries[0].URI)
			originalURI, err := url.Parse(mediaURI)
			if err != nil {
				logger.Errorf("Failed to parse original URI %s for stream %s: %v", mediaURI, s.MediaName(), err)
				return contenttype.MediaType{}, err
			}
			uri.Scheme = originalURI.Scheme
			uri.Host = originalURI.Host
			if !strings.HasPrefix(uri.Path, "/") {
				uri.Path = path.Join(path.Dir(originalURI.Path), uri.Path)
			}
			logger.Debugf("Resolved relative URL for stream %s to: %s", s.MediaName(), uri.String())
		}

		return s.verify(uri.String())
	}

	logger.Debugf("Media verification successful for stream %s: content_type=%s", s.MediaName(), ct.String())
	return ct, nil
}

func (s *BaseStreamSource) verifyWithDiags(mediaURI string, diag *StreamSourceDiag) {
	logger.Debugf("Verifying media URI with diagnostics for stream %s: %s", s.MediaName(), mediaURI)

	s.mux.RLock()
	body, status, ct, err := s.conn.Get("GET", mediaURI)
	s.mux.RUnlock()
	if err != nil {
		logger.Errorf("Failed to fetch media from %s for stream %s: status=%d, error=%v",
			mediaURI, s.MediaName(), status, err)
		diag.Diagnostics = append(diag.Diagnostics, HttpDiags{
			Url:    mediaURI,
			Status: status,
			Error:  err.Error(),
		})
		return
	}

	logger.Debugf("Retrieved media for diagnostic of stream %s: status=%d, content_type=%s, body_size=%d bytes",
		s.MediaName(), status, ct.String(), len(body))

	httpDiag := HttpDiags{
		Url:       mediaURI,
		Body:      string(body),
		MediaType: ct.String(),
		Status:    status,
	}

	if !ct.MatchesAny(supportedMediaTypes...) {
		logger.Warnf("Unsupported content type %s for stream %s at %s", ct.String(), s.MediaName(), mediaURI)
		httpDiag.Error = "invalid content type"
		diag.Diagnostics = append(diag.Diagnostics, httpDiag)
		return
	}

	if ct.Subtype == "vnd.apple.mpegurl" || ct.Subtype == "x-mpegurl" {
		logger.Debugf("Processing M3U playlist for diagnostic of stream %s", s.MediaName())
		m3uPlaylist, err := m3uparser.DecodeFromReader(bytes.NewReader(body))
		if err != nil {
			logger.Errorf("Failed to parse M3U playlist for stream %s: %v", s.MediaName(), err)
			httpDiag.Error = err.Error()
			diag.Diagnostics = append(diag.Diagnostics, httpDiag)
			return
		}

		if len(m3uPlaylist.Entries) == 0 {
			logger.Warnf("Empty M3U playlist for stream %s at %s", s.MediaName(), mediaURI)
			httpDiag.Error = "empty playlist"
			diag.Diagnostics = append(diag.Diagnostics, httpDiag)
			return
		}

		logger.Debugf("M3U playlist for stream %s contains %d entries, following first entry for diagnostic",
			s.MediaName(), len(m3uPlaylist.Entries))

		uri, _ := url.Parse(m3uPlaylist.Entries[0].URI)

		if uri.Scheme == "" {
			logger.Debugf("Resolving relative URL for diagnostic of stream %s: %s", s.MediaName(), m3uPlaylist.Entries[0].URI)
			originalURI, err := url.Parse(mediaURI)
			if err != nil {
				logger.Errorf("Failed to parse original URI %s for stream %s: %v", mediaURI, s.MediaName(), err)
				httpDiag.Error = err.Error()
				diag.Diagnostics = append(diag.Diagnostics, httpDiag)
				return
			}
			uri.Scheme = originalURI.Scheme
			uri.Host = originalURI.Host
			if !strings.HasPrefix(uri.Path, "/") {
				uri.Path = path.Join(path.Dir(originalURI.Path), uri.Path)
			}
			logger.Debugf("Resolved relative URL for diagnostic of stream %s to: %s", s.MediaName(), uri.String())
		}

		s.verifyWithDiags(uri.String(), diag)
	}

	diag.Diagnostics = append(diag.Diagnostics, httpDiag)
	diag.Active = status == http.StatusOK

	logger.Debugf("Diagnostic verification completed for stream %s: active=%t, final_status=%d",
		s.MediaName(), diag.Active, status)
}
