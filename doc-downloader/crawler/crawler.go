package crawler

import (
	"net/url"
	"sync"

	"github.com/ellenzogalla/doc-downloader.git/downloader"
	"github.com/ellenzogalla/doc-downloader.git/utils"

	"github.com/gocolly/colly"
	"github.com/playwright-community/playwright-go"
)

// Crawler manages the web crawling process.
type Crawler struct {
	targetURL  string
	outputDir  string
	wg         *sync.WaitGroup
	browsers   chan *playwright.Browser
	colly      *colly.Collector
	visited    map[string]bool
	visitedMux sync.Mutex
}

// New creates a new Crawler instance.
func New(targetURL, outputDir string, wg *sync.WaitGroup, browsers chan *playwright.Browser) *Crawler {
	c := &Crawler{
		targetURL:  targetURL,
		outputDir:  outputDir,
		wg:         wg,
		browsers:   browsers,
		visited:    make(map[string]bool),
		visitedMux: sync.Mutex{},
	}

	// Initialize Colly
	c.colly = colly.NewCollector(
		colly.Async(true), // Enable asynchronous requests
	)

	// Limit parallelism
	c.colly.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2, // Adjust as needed
	})

	// Set up request callback (for downloading)
	c.colly.OnRequest(func(r *colly.Request) {
		c.visitedMux.Lock()
		if c.visited[r.URL.String()] {
			c.visitedMux.Unlock()
			r.Abort() // Skip if already visited
			return
		}
		c.visited[r.URL.String()] = true
		c.visitedMux.Unlock()

		c.wg.Add(1)
		browser := <-c.browsers // Get a browser from the pool
		go func(urlStr string) {
			downloader.DownloadHTMLWithInlineStyles(urlStr, c.outputDir, browser, c.wg)
			c.browsers <- browser // Return browser to the pool
		}(r.URL.String())
	})

	// Set up HTML parsing callback (for link discovery)
	c.colly.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if c.shouldVisit(link) {
			e.Request.Visit(link)
		}
	})

	return c
}

// Crawl starts the web crawling process.
func (c *Crawler) Crawl() {
	c.colly.Visit(c.targetURL)
	c.colly.Wait() // Wait for Colly to finish its tasks
}

// shouldVisit checks if a URL should be visited.
func (c *Crawler) shouldVisit(link string) bool {
	c.visitedMux.Lock()
	defer c.visitedMux.Unlock()

	if c.visited[link] {
		return false // Already visited
	}

	parsedLink, err := url.Parse(link)
	if err != nil {
		return false // Invalid URL
	}

	// Check if the link is within the same domain as the target URL
	return parsedLink.Hostname() == utils.GetBaseHostname(c.targetURL)
}
