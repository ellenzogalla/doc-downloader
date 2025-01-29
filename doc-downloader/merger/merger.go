package merger

import (
	"io/ioutil"
	"path/filepath"
	"sort"
)

// MergeHTMLFiles merges all HTML files in a directory into a single HTML file.
func MergeHTMLFiles(dir, outputFilename string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Sort files for consistent merging (you might need more sophisticated sorting)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var mergedContent string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			filePath := filepath.Join(dir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			mergedContent += string(content) + "\n"
		}
	}

	outputPath := filepath.Join(dir, outputFilename)
	err = ioutil.WriteFile(outputPath, []byte(mergedContent), 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
