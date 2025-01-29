package converter

import (
	"fmt"
	"log"
	"os/exec"
)

// ConvertToPDF converts an HTML file to a PDF using wkhtmltopdf.
func ConvertToPDF(htmlFilePath, pdfFilePath string) error {
	// Ensure wkhtmltopdf is installed
	_, err := exec.LookPath("wkhtmltopdf")
	if err != nil {
		return fmt.Errorf("wkhtmltopdf not found. Please install it: %v", err)
	}

	// Construct the command
	cmd := exec.Command(
		"wkhtmltopdf",
		"--enable-local-file-access", // Allow access to local files
		"--page-size", "A4",          // Set page size to A4
		"--margin-top", "25mm", // Adjust margins as needed
		"--margin-bottom", "25mm",
		"--margin-left", "20mm",
		"--margin-right", "20mm",
		htmlFilePath,
		pdfFilePath,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("wkhtmltopdf output: %s", output) // Log wkhtmltopdf output
		return fmt.Errorf("failed to convert HTML to PDF: %v", err)
	}

	return nil
}
