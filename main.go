package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	urlsFile := flag.String("urls", "", "Path to file containing URLs (required)")
	outDir := flag.String("out", "./downloads", "Directory to save files")
	workers := flag.Int("workers", 4, "Max concurrent downloads")
	flag.Parse()

	if *urlsFile == "" {
		fmt.Fprintln(os.Stderr, "error: -urls flag is required")
		flag.Usage()
		os.Exit(1)
	}

	urls, err := readUrls(*urlsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading urls file:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "error creating output dir:", err)
		os.Exit(1)
	}

	start := time.Now()
	results := runDownloads(urls, *outDir, *workers)

	var succeeded, failed int
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("  ✗ %-30s failed: %v\n", r.Filename, r.Err)
			failed++
		} else {
			fmt.Printf("  ✓ %-30s (%.1f KB, %.2fs)\n",
				r.Filename, float64(r.Bytes)/1024, r.Duration.Seconds())
			succeeded++
		}
	}

	fmt.Printf("\nDone. %d succeeded, %d failed. Total time: %.2fs\n",
		succeeded, failed, time.Since(start).Seconds())
}
