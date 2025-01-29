package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
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

// ExtractLinks finds all the links within an HTML document and returns
// absolute URLs that are within the same domain as the base URL.
func ExtractLinks(htmlContent []byte, baseURL string) []string {
	var links []string
	base, _ := url.Parse(baseURL)

	tokenizer := html.NewTokenizer(strings.NewReader(string(htmlContent)))
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					linkURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					absURL := base.ResolveReference(linkURL)

					// Only add links from the same domain
					if absURL.Hostname() == base.Hostname() {
						links = append(links, absURL.String())
					}
				}
			}
		}
	}

	return links
}

// GetFilePath generates a file path based on the URL and output directory.
func GetFilePath(outputDir, rawURL, extension string) string {
	parsedURL, _ := url.Parse(rawURL)
	path := parsedURL.Path

	// Remove leading slash
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Replace slashes with underscores
	path = strings.ReplaceAll(path, "/", "_")

	// Remove trailing underscore if it's from a directory index
	if strings.HasSuffix(path, "_") {
		path = path[:len(path)-1]
	}

	// Add extension if not already present
	if !strings.HasSuffix(path, extension) {
		path += extension
	}

	// Create directory structure if it doesn't exist
	dir := filepath.Dir(path)
	fullDirPath := filepath.Join(outputDir, dir)
	os.MkdirAll(fullDirPath, 0755)

	return filepath.Join(outputDir, path)
}
