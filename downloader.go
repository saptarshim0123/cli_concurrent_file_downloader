package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Result struct {
	URL      string
	Filename string
	Bytes    int64
	Duration time.Duration
	Err      error
}

func download(url, outDir string) Result {
	start := time.Now()

	filename := filepath.Base(url)
	if filename == "." || filename == "/" {
		filename = "index.html"
	}
	destPath := filepath.Join(outDir, filename)

	resp, err := http.Get(url)
	if err != nil {
		return Result{URL: url, Filename: filename, Err: fmt.Errorf("request failed: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{
			URL:      url,
			Filename: filename,
			Err:      fmt.Errorf("bad status: %s", resp.Status),
		}
	}

	out, err := os.Create(destPath)
	if err != nil {
		return Result{URL: url, Filename: filename, Err: fmt.Errorf("create file failed: %w", err)}
	}
	defer out.Close()

	// stream response body -> file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return Result{URL: url, Filename: filename, Err: fmt.Errorf("write failed: %w", err)}
	}

	return Result{
		URL:      url,
		Filename: filename,
		Bytes:    n,
		Duration: time.Since(start),
	}
}

func runDownloads(urls []string, outDir string, numWorkers int) []Result {
	total := len(urls)

	// buffered channels
	jobs := make(chan string, total)
	results := make(chan Result, total)

	// launch N workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for url := range jobs { // range over channel blocks until job arrives
				fmt.Printf("  → Downloading: %s\n", url)
				results <- download(url, outDir) // send result when done
			}
		}()
	}

	// feed all URLs into jobs channel
	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // signal workers: no more jobs coming

	// collect all results
	all := make([]Result, 0, total)
	for i := 0; i < total; i++ {
		all = append(all, <-results)
	}
	return all
}
