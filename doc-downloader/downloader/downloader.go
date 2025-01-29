package downloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/ellenzogalla/doc-downloader.git/utils"

	"github.com/playwright-community/playwright-go"
)

// DownloadHTMLWithInlineStyles downloads the page using Playwright, waits for it to be fully rendered,
// and then saves the HTML content with inline styles and Base64 encoded images.
func DownloadHTMLWithInlineStyles(url, outputDir string, browser *playwright.Browser, wg *sync.WaitGroup) {
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

	// Evaluate JavaScript to inline styles and convert images to Base64
	_, err = page.Evaluate(`
		() => {
			// Inline styles
			for (const style of document.querySelectorAll('style')) {
				style.textContent = style.textContent + ''; // Forces re-evaluation of styles
			}
			for (const el of document.querySelectorAll('*')) {
				if (el.style) {
					const computedStyle = getComputedStyle(el);
					for (let i = 0; i < computedStyle.length; i++) {
						const property = computedStyle[i];
						el.style[property] = computedStyle.getPropertyValue(property);
					}
				}
			}

			// Convert images to Base64
			const images = document.querySelectorAll('img');
			for (const img of images) {
				if (img.src.startsWith('data:')) continue; // Already Base64
				const canvas = document.createElement('canvas');
				const ctx = canvas.getContext('2d');
				canvas.width = img.width;
				canvas.height = img.height;
				ctx.drawImage(img, 0, 0);
				try {
					img.src = canvas.toDataURL();
				} catch (e) {
					console.error('Failed to convert image to Base64:', e);
				}
			}
		}
	`)
	if err != nil {
		fmt.Println("Failed to inline styles and convert images to Base64:", err)
		return
	}

	// Get the HTML content after inlining styles and converting images
	htmlContent, err := page.Content()
	if err != nil {
		fmt.Println("Failed to get page content:", err)
		return
	}

	// Save HTML with inline styles and Base64 images
	htmlFilePath := utils.GetFilePath(outputDir, url, ".html")
	err = os.WriteFile(htmlFilePath, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Println("Failed to save HTML:", err)
		return
	}
	fmt.Println("Downloaded (HTML with inline styles):", url)
}

// CombineHTMLFiles reads all HTML files in the output directory and combines their content.
func CombineHTMLFiles(outputDir string) (string, error) {
	var combinedHTML string

	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to read output directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" && file.Name() != "combined.html" {
			filePath := filepath.Join(outputDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read HTML file %s: %v", file.Name(), err)
			}
			combinedHTML += string(content) + "\n"
		}
	}

	return combinedHTML, nil
}

// ConvertToPDF converts the HTML to a PDF using Playwright
func ConvertToPDF(htmlFilePath, pdfFilePath string, browser *playwright.Browser) error {
	page, err := (*browser).NewPage()
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err)
	}
	defer page.Close()

	absHTMLFilePath, err := filepath.Abs(htmlFilePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for HTML file: %v", err)
	}

	// Use file:// protocol to open the local HTML file
	if _, err = page.Goto("file://"+absHTMLFilePath, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return fmt.Errorf("failed to open HTML file: %v", err)
	}

	_, err = page.PDF(playwright.PagePdfOptions{
		Path:   playwright.String(pdfFilePath),
		Format: playwright.String("A4"),
	})
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %v", err)
	}

	fmt.Println("Converted to PDF:", pdfFilePath)
	return nil
}
