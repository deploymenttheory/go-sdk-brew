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

	// Returns analytics scoped to homebrew-core formulae only.
	// Response shape: formulae map (formula name → []CoreFormulaItem) rather than a ranked items array.
	category := analytics.CategoryInstall
	period := analytics.Period30d

	result, _, err := client.Analytics.ListCoreByCategory(context.Background(), category, period)
	if err != nil {
		log.Fatalf("failed to list core analytics: %v", err)
	}

	fmt.Printf("Category:    %s\n", result.Category)
	fmt.Printf("Period:      %s – %s\n", result.StartDate, result.EndDate)
	fmt.Printf("Total items: %d\n", result.TotalItems)
	fmt.Printf("Total count: %d\n\n", result.TotalCount)

	// Print a sample of entries.
	printed := 0
	for name, entries := range result.Formulae {
		if printed >= 5 {
			break
		}
		entryJSON, _ := json.Marshal(entries)
		fmt.Printf("  %-30s %s\n", name, entryJSON)
		printed++
	}
}
