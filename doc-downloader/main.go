package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ellenzogalla/doc-downloader.git/crawler"
	"github.com/ellenzogalla/doc-downloader.git/merger"
	"github.com/ellenzogalla/doc-downloader.git/utils"
	"github.com/playwright-community/playwright-go"
)

func main() {
	// Command-line flags
	targetURL := flag.String("url", "", "The base URL of the documentation website")
	outputDir := flag.String("out", "output", "The directory to save downloaded files")
	numWorkers := flag.Int("workers", 4, "Number of worker processes (for Playwright instances)")
	flag.Parse()

	if *targetURL == "" {
		log.Fatal("Error: Please provide the target URL using the -url flag.")
	}

	// Create output directory
	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Could not start Playwright: %v", err)
	}
	defer pw.Stop()

	// Browser pool for concurrent tasks
	browsers := make(chan *playwright.Browser, *numWorkers)
	for i := 0; i < *numWorkers; i++ {
		browser, err := pw.Chromium.Launch()
		if err != nil {
			log.Fatalf("Could not launch browser: %v", err)
		}
		browsers <- &browser
	}

	// Crawl, download, and convert each page to PDF
	var wg sync.WaitGroup
	c := crawler.New(*targetURL, *outputDir, &wg, browsers)
	c.Crawl()

	// Wait for all downloads and conversions to complete
	wg.Wait()

	// Merge PDFs into a single file
	pdfFiles, err := utils.ListPDFFiles(*outputDir)
	if err != nil {
		log.Fatalf("Error listing PDF files: %v", err)
	}

	mergedPDFPath := filepath.Join(*outputDir, "final.pdf")
	err = merger.MergePDFs(pdfFiles, mergedPDFPath)
	if err != nil {
		log.Fatalf("Error merging PDFs: %v", err)
	}

	fmt.Println("Documentation downloaded and converted to a single PDF:", mergedPDFPath)
}
