package upstream

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/a13labs/m3uproxy/pkg/logger"
	"github.com/elnormous/contenttype"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func NewUpstreamConnection(headers map[string]string, proxy string, timeout int) *UpstreamConnection {
	logger.Debugf("Creating new upstream connection with timeout: %ds, proxy: %s", timeout, proxy)

	var dial fasthttp.DialFunc

	if proxy != "" {
		logger.Debugf("Configuring proxy connection: %s", proxy)
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			logger.Errorf("Invalid proxy URL '%s': %s", proxy, err)
			panic(fmt.Sprintf("Invalid proxy URL: %s", err))
		}

		dial = fasthttpproxy.FasthttpHTTPDialer(proxyURL.Host)
		logger.Infof("Upstream connection configured with proxy: %s", proxyURL.Host)
	} else {
		logger.Debug("Upstream connection configured without proxy")
		dial = fasthttp.Dial
	}

	logger.Debugf("Upstream connection created with %d custom headers", len(headers))
	return &UpstreamConnection{
		client: &fasthttp.Client{
			ReadTimeout: time.Duration(timeout) * time.Second,
			Dial:        dial,
		},
		headers: headers,
	}
}

func (u *UpstreamConnection) Check(method string, uri string) (string, contenttype.MediaType, error) {
	logger.Debugf("Starting upstream check for %s %s", method, uri)

	const maxRedirects = 10
	currentURL := uri
	var resp *fasthttp.Response
	redirectCount := 0

	for i := 0; i < maxRedirects; i++ {
		logger.Debugf("Attempt %d/%d: Checking URL %s", i+1, maxRedirects, currentURL)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(currentURL)
		req.Header.SetMethod(method)

		for key, value := range u.headers {
			req.Header.Set(key, value)
		}

		resp = fasthttp.AcquireResponse()
		err := u.client.Do(req, resp)
		if err != nil {
			logger.Errorf("HTTP request failed for %s: %v", currentURL, err)
			fasthttp.ReleaseResponse(resp)
			return "", contenttype.MediaType{}, err
		}

		statusCode := resp.StatusCode()
		ct := contenttype.NewMediaType(string(resp.Header.ContentType()))
		logger.Debugf("Received response: status=%d, content-type=%s", statusCode, ct.String())

		if statusCode/100 == 3 { // Handle redirects (3xx status codes)
			redirectCount++
			logger.Debugf("Redirect %d: status=%d", redirectCount, statusCode)

			location := resp.Header.Peek("Location")
			if location == nil {
				logger.Errorf("Redirect response missing Location header for %s", currentURL)
				fasthttp.ReleaseResponse(resp)
				return "", ct, fmt.Errorf("redirect response missing Location header")
			}

			// Resolve the new URL relative to the current URL
			newURL := string(location)
			logger.Debugf("Redirect location: %s", newURL)

			if !strings.HasPrefix(newURL, "http") {
				baseURL, err := url.Parse(currentURL)
				if err != nil {
					logger.Errorf("Failed to parse base URL '%s': %v", currentURL, err)
					fasthttp.ReleaseResponse(resp)
					return "", ct, fmt.Errorf("failed to parse base URL: %w", err)
				}
				relativeURL, err := url.Parse(newURL)
				if err != nil {
					logger.Errorf("Failed to parse relative URL '%s': %v", newURL, err)
					fasthttp.ReleaseResponse(resp)
					return "", ct, fmt.Errorf("failed to parse relative URL: %w", err)
				}
				currentURL = baseURL.ResolveReference(relativeURL).String()
				logger.Debugf("Resolved relative URL to: %s", currentURL)
			} else {
				currentURL = newURL
			}

			// Release the response and continue to the next redirect
			fasthttp.ReleaseResponse(resp)
			continue
		}

		logger.Infof("Upstream check completed for %s: final_url=%s, status=%d, content_type=%s, redirects=%d",
			uri, currentURL, statusCode, ct.String(), redirectCount)
		return currentURL, ct, nil
	}

	// Exceeded maximum redirects
	logger.Errorf("Too many redirects (%d) for %s, last URL: %s", maxRedirects, uri, currentURL)
	if resp != nil {
		fasthttp.ReleaseResponse(resp)
	}
	return "", contenttype.MediaType{}, fmt.Errorf("too many redirects")
}

func (u *UpstreamConnection) Get(method, URI string) ([]byte, int, contenttype.MediaType, error) {
	logger.Debugf("Starting upstream GET request: %s %s", method, URI)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(URI)
	req.Header.SetMethod(method)

	for key, value := range u.headers {
		req.Header.Set(key, value)
	}

	logger.Debugf("Sending request to %s with %d custom headers", URI, len(u.headers))

	resp := fasthttp.AcquireResponse()
	err := u.client.Do(req, resp)
	if err != nil {
		logger.Errorf("HTTP request failed for %s %s: %v", method, URI, err)
		fasthttp.ReleaseResponse(resp)
		return nil, -1, contenttype.MediaType{}, err
	}

	statusCode := resp.StatusCode()
	ct := contenttype.NewMediaType(string(resp.Header.ContentType()))
	body := resp.Body()
	bodySize := len(body)

	logger.Debugf("Received response from %s: status=%d, content_type=%s, body_size=%d bytes",
		URI, statusCode, ct.String(), bodySize)

	if statusCode/100 != 2 {
		logger.Warnf("Non-2xx response from %s: status=%d", URI, statusCode)
		fasthttp.ReleaseResponse(resp)
		return nil, statusCode, ct, fmt.Errorf("http response code (%d)", statusCode)
	}

	if statusCode == fasthttp.StatusNoContent {
		logger.Infof("No content response from %s", URI)
		fasthttp.ReleaseResponse(resp)
		return nil, statusCode, ct, errors.New("no content")
	}

	// Detect gzip either via Content-Encoding header or magic bytes and decompress
	enc := strings.ToLower(string(resp.Header.Peek("Content-Encoding")))
	if strings.Contains(enc, "gzip") || (len(body) >= 2 && body[0] == 0x1f && body[1] == 0x8b) {
		logger.Debugf("Detected gzip-compressed body for %s, attempting to decompress", URI)
		gr, gzerr := gzip.NewReader(bytes.NewReader(body))
		if gzerr == nil {
			decomp, reaerr := io.ReadAll(gr)
			_ = gr.Close()
			if reaerr == nil {
				body = decomp
				bodySize = len(body)
				logger.Debugf("Decompression succeeded for %s, new body_size=%d", URI, bodySize)
			} else {
				logger.Warnf("Failed to read gzip body for %s: %v", URI, reaerr)
			}
		} else {
			logger.Warnf("Failed to create gzip reader for %s: %v", URI, gzerr)
		}
	}

	dst := make([]byte, len(body))
	copy(dst, body)

	// release response now that we copied the body to our buffer
	fasthttp.ReleaseResponse(resp)

	logger.Infof("Successfully retrieved %d bytes from %s %s", len(dst), method, URI)
	return dst, statusCode, ct, nil
}
