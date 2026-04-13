// Recipe: Top Formulae Report
//
// Fetches the top N most-installed formulae from the 30-day install analytics,
// then enriches each entry with full formula metadata (description, license,
// homepage, dependencies) in a single formatted report.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/formulae"
)

const topN = 10

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Fetch install analytics for 30 days.
	stats, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
	if err != nil {
		log.Fatalf("failed to fetch analytics: %v", err)
	}

	items := stats.Items
	if len(items) > topN {
		items = items[:topN]
	}

	fmt.Printf("╔══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║         Top %d Most-Installed Homebrew Formulae (30 days)            ║\n", topN)
	fmt.Printf("║         Period: %s → %s                              ║\n", stats.StartDate, stats.EndDate)
	fmt.Printf("╚══════════════════════════════════════════════════════════════════════╝\n\n")

	// Step 2: Enrich each entry with full formula metadata.
	for _, item := range items {
		formula, _, err := client.Formulae.GetByName(ctx, item.Formula)
		if err != nil {
			fmt.Printf("#%-3d %-25s %s installs (%s%%)\n",
				item.Number, item.Formula, item.Count, item.Percent)
			fmt.Printf("     [metadata unavailable: %v]\n\n", err)
			continue
		}

		printFormulaRow(item, formula)
	}

	fmt.Printf("\nTotal formulae in catalog: %d | Total installs recorded: %s\n",
		stats.TotalItems, formatCount(stats.TotalCount))
}

func printFormulaRow(item analytics.AnalyticsItem, f *formulae.ResourceFormula) {
	fmt.Printf("#%-3d %-25s %s installs (%s%%)\n",
		item.Number, item.Formula, item.Count, item.Percent)
	fmt.Printf("     Description:  %s\n", f.Desc)
	fmt.Printf("     Version:      %s\n", f.Versions.Stable)
	fmt.Printf("     License:      %s\n", f.License)
	fmt.Printf("     Homepage:     %s\n", f.Homepage)
	if len(f.Dependencies) > 0 {
		fmt.Printf("     Dependencies: %s\n", strings.Join(f.Dependencies, ", "))
	}
	fmt.Println()
}

func formatCount(n int) string {
	// Simple thousands-separator formatting.
	s := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
