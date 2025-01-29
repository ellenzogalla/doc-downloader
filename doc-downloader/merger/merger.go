package merger

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// MergePDFs merges multiple PDF files into a single PDF.
func MergePDFs(pdfFiles []string, outputFilename string) error {
	config := pdfcpu.NewDefaultConfiguration()
	return api.MergeCreateFile(pdfFiles, outputFilename, config)
}
