package main

import (
	"context"
	"fmt"
	"log"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
)

func main() {
	// Create a client using the default Homebrew API base URL.
	// No authentication required — the Homebrew API is public.
	client, err := brew.NewClient("",
		brew.WithTimeout(30_000_000_000), // 30s
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// List all formulae.
	formulae, _, err := client.Formulae.List(ctx)
	if err != nil {
		log.Fatalf("failed to list formulae: %v", err)
	}
	fmt.Printf("Total formulae: %d\n", len(formulae))
	if len(formulae) > 0 {
		fmt.Printf("First formula: %s (%s)\n", formulae[0].Name, formulae[0].Versions.Stable)
	}

	// Get a single formula by name.
	wget, _, err := client.Formulae.GetByName(ctx, "wget")
	if err != nil {
		log.Fatalf("failed to get wget: %v", err)
	}
	fmt.Printf("\nwget v%s: %s\n", wget.Versions.Stable, wget.Desc)
	fmt.Printf("  Homepage: %s\n", wget.Homepage)
	fmt.Printf("  License: %s\n", wget.License)
	fmt.Printf("  Dependencies: %v\n", wget.Dependencies)

	// List all casks.
	casks, _, err := client.Casks.List(ctx)
	if err != nil {
		log.Fatalf("failed to list casks: %v", err)
	}
	fmt.Printf("\nTotal casks: %d\n", len(casks))

	// Get a single cask by name.
	iterm2, _, err := client.Casks.GetByName(ctx, "iterm2")
	if err != nil {
		log.Fatalf("failed to get iterm2: %v", err)
	}
	fmt.Printf("\niterm2 v%s: %s\n", iterm2.Version, iterm2.Desc)

	// Get install analytics for the last 30 days.
	stats, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
	if err != nil {
		log.Fatalf("failed to get analytics: %v", err)
	}
	fmt.Printf("\nTop 3 most installed formulae (30d):\n")
	for i, item := range stats.Items {
		if i >= 3 {
			break
		}
		fmt.Printf("  %d. %s (%s installs, %s)\n", item.Number, item.Formula, item.Count, item.Percent)
	}
}
