package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
)

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Returns cask install analytics scoped to homebrew-cask.
	// Response shape: formulae map (cask token → []CaskInstallItem) rather than a ranked items array.
	period := analytics.Period30d

	result, _, err := client.Analytics.ListCaskInstalls(context.Background(), period)
	if err != nil {
		log.Fatalf("failed to list cask install analytics: %v", err)
	}

	fmt.Printf("Category:    %s\n", result.Category)
	fmt.Printf("Period:      %s – %s\n", result.StartDate, result.EndDate)
	fmt.Printf("Total items: %d\n", result.TotalItems)
	fmt.Printf("Total count: %d\n\n", result.TotalCount)

	// Print a sample of entries.
	printed := 0
	for token, entries := range result.Formulae {
		if printed >= 5 {
			break
		}
		entryJSON, _ := json.Marshal(entries)
		fmt.Printf("  %-30s %s\n", token, entryJSON)
		printed++
	}
}
