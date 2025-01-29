package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ellenzogalla/doc-downloader.git/converter"
	"github.com/ellenzogalla/doc-downloader.git/downloader"
	"github.com/ellenzogalla/doc-downloader.git/queue"
	"github.com/ellenzogalla/doc-downloader.git/utils"
)

func main() {
	// Command-line flags
	targetURL := flag.String("url", "", "The base URL of the documentation website")
	outputDir := flag.String("out", "output", "The directory to save downloaded files")
	numWorkers := flag.Int("workers", 4, "Number of worker processes")
	flag.Parse()

	if *targetURL == "" {
		log.Fatal("Error: Please provide the target URL using the -url flag.")
	}

	// Create output directory if it doesn't exist
	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	// Task queue and synchronization
	taskQueue := queue.NewTaskQueue()
	var wg sync.WaitGroup

	// Start worker processes
	for i := 0; i < *numWorkers; i++ {
		wg.Add(1)
		go worker(*outputDir, taskQueue, &wg)
	}

	// Seed the queue with the initial URL
	baseURL, err := utils.NormalizeBaseURL(*targetURL)
	if err != nil {
		log.Fatal("Error normalizing base URL:", err)
	}
	taskQueue.Enqueue(queue.Task{URL: baseURL, Type: queue.TaskTypeDownload})

	// Crawl the website
	visited := make(map[string]bool)
	for {
		task, ok := taskQueue.Dequeue()
		if !ok {
			if taskQueue.IsEmpty() {
				break // Queue is empty, crawling finished
			}
			time.Sleep(100 * time.Millisecond) // Wait a bit for new tasks
			continue
		}

		if visited[task.URL] {
			continue // Already processed
		}
		visited[task.URL] = true

		if task.Type == queue.TaskTypeDownload {
			page, err := downloader.Download(task.URL)
			if err != nil {
				log.Printf("Error downloading %s: %v", task.URL, err)
				continue
			}

			// Save the HTML to the output directory
			filePath := utils.GetFilePath(*outputDir, task.URL, ".html")
			err = downloader.Save(page, filePath)
			if err != nil {
				log.Printf("Error saving HTML for %s: %v", task.URL, err)
				continue
			}
			fmt.Println("Downloaded:", task.URL)

			// Enqueue PDF conversion task
			taskQueue.Enqueue(queue.Task{URL: task.URL, Type: queue.TaskTypeConvert, FilePath: filePath})

			// Find and enqueue new links
			links := utils.ExtractLinks(page, baseURL)
			for _, link := range links {
				if !visited[link] {
					taskQueue.Enqueue(queue.Task{URL: link, Type: queue.TaskTypeDownload})
				}
			}
		}
	}

	// Wait for workers to finish
	wg.Wait()
	fmt.Println("Documentation download and conversion complete.")
}

// Worker function for processing tasks
func worker(outputDir string, taskQueue *queue.TaskQueue, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := taskQueue.Dequeue()
		if !ok {
			break // Queue is closed
		}

		if task.Type == queue.TaskTypeConvert {
			pdfFilePath := utils.GetFilePath(outputDir, task.URL, ".pdf")
			err := converter.ConvertToPDF(task.FilePath, pdfFilePath)
			if err != nil {
				log.Printf("Error converting to PDF %s: %v", task.URL, err)
			} else {
				fmt.Println("Converted to PDF:", task.URL)
			}
		}
	}
}
