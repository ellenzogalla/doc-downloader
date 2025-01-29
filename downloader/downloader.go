package downloader

import (
	"io/ioutil"
	"net/http"
)

// Download fetches the content of a URL.
func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Save writes the content to a file.
func Save(content []byte, filePath string) error {
	return ioutil.WriteFile(filePath, content, 0644)
}
