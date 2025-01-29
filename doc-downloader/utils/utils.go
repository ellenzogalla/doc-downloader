package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// NormalizeBaseURL ensures the base URL has a trailing slash.
func NormalizeBaseURL(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return u.String(), nil
}

// GetFilePath generates a file path based on the URL and output directory.
func GetFilePath(outputDir, rawURL, extension string) string {
	parsedURL, _ := url.Parse(rawURL)
	path := parsedURL.Path

	// Remove leading slash
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Replace slashes with underscores, preserving directory structure
	path = strings.ReplaceAll(path, "/", string(os.PathSeparator))

	// Add extension if not already present
	if !strings.HasSuffix(path, extension) {
		path += extension
	}

	// Create directory structure if it doesn't exist
	fullPath := filepath.Join(outputDir, path)
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, 0755)

	return fullPath
}

// GetBaseHostname extracts the base hostname from a URL.
func GetBaseHostname(rawURL string) string {
	parsedURL, _ := url.Parse(rawURL)
	return parsedURL.Hostname()
}
