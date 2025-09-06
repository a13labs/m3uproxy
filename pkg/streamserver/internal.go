package streamserver

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"
)

func loadContent(filePath string) (string, error) {
	var reader io.Reader

	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		resp, err := http.Get(filePath)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		reader = resp.Body
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()
		reader = file
	}

	// Read first few bytes to check for gzip magic number
	buf := make([]byte, 512)
	n, err := io.ReadFull(reader, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", err
	}
	data := buf[:n]

	// Check for gzip magic number
	isGzip := n >= 2 && data[0] == 0x1f && data[1] == 0x8b

	var content []byte
	if isGzip {
		gzReader, err := gzip.NewReader(io.MultiReader(bytes.NewReader(data), reader))
		if err != nil {
			return "", err
		}
		defer gzReader.Close()
		content, err = io.ReadAll(gzReader)
		if err != nil {
			return "", err
		}
	} else {
		content = data
		rest, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		content = append(content, rest...)
	}

	return string(content), nil
}
