package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readUrls(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
        return nil, fmt.Errorf("could not open urls file: %w", err)
    }
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	return urls, scanner.Err()
}