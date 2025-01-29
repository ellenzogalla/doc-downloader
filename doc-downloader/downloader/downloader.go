package downloader

import (
	"fmt"
	"os"
	"sync"

	"github.com/ellenzogalla/doc-downloader.git/utils"
	"github.com/playwright-community/playwright-go"
)

// DownloadWithPlaywright downloads the page using Playwright, waits for it to be fully rendered,
// and then saves the HTML content and converts it to PDF.
func DownloadWithPlaywright(url, outputDir string, browser *playwright.Browser, wg *sync.WaitGroup) {
	defer wg.Done()

	page, err := (*browser).NewPage()
	if err != nil {
		fmt.Println("Failed to create page:", err)
		return
	}
	defer page.Close()

	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		fmt.Println("Failed to goto:", err)
		return
	}

	htmlContent, err := page.Content()
	if err != nil {
		fmt.Println("Failed to get page content:", err)
		return
	}

	// Save HTML
	htmlFilePath := utils.GetFilePath(outputDir, url, ".html")
	err = os.WriteFile(htmlFilePath, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("Failed to save HTML:", err)
		return
	}
	fmt.Println("Downloaded (HTML):", url)

	// Convert to PDF (using Playwright)
	pdfFilePath := utils.GetFilePath(outputDir, url, ".pdf")
	_, err = page.PDF(playwright.PagePdfOptions{
		Path:   playwright.String(pdfFilePath),
		Format: playwright.String("A4"),
	})
	if err != nil {
		fmt.Println("Failed to generate PDF:", err)
		return
	}
	fmt.Println("Converted to PDF:", url)
}
