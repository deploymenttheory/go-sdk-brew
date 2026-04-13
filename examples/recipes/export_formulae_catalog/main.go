// Recipe: Export Formulae Catalog
//
// Fetches the complete Homebrew formula catalog from /api/formula.json and
// writes it to a local JSON file. Useful for offline analysis, caching, or
// feeding into downstream tooling without repeated API calls.
//
// Usage:
//
//	go run main.go                          # writes formulae.json in the current directory
//	go run main.go -output /tmp/brew.json   # writes to a custom path
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
)

func main() {
	output := flag.String("output", "formulae.json", "path to write the exported JSON file")
	flag.Parse()

	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	fmt.Println("Fetching formula catalog from https://formulae.brew.sh …")
	start := time.Now()

	formulae, resp, err := client.Formulae.List(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch formulae: %v", err)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	fmt.Printf("Received %d formulae in %s (HTTP %d)\n", len(formulae), elapsed, resp.StatusCode())

	// Marshal with indentation for human readability.
	data, err := json.MarshalIndent(formulae, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal formulae: %v", err)
	}

	if err := os.WriteFile(*output, data, 0o644); err != nil {
		log.Fatalf("failed to write %s: %v", *output, err)
	}

	info, err := os.Stat(*output)
	if err != nil {
		log.Fatalf("failed to stat output file: %v", err)
	}

	fmt.Printf("Exported %d formulae → %s (%.1f MB)\n",
		len(formulae), *output, float64(info.Size())/1024/1024)
}
