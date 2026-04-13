// Recipe: Trending Packages
//
// Identifies which Homebrew formulae are trending upward by comparing their
// install count growth rate across three periods (30d, 90d, 365d).
// For each formula in the 30d top 50, it computes an implied monthly rate for
// the 90d and 365d windows and flags those growing faster than average.
package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
)

type formulaGrowth struct {
	name       string
	count30d   float64
	rate90d    float64 // avg monthly installs over 90d window
	rate365d   float64 // avg monthly installs over 365d window
	momentum   float64 // count30d / rate90d — >1 means accelerating
}

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Fetch install analytics for all three periods.
	stats30, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
	if err != nil {
		log.Fatalf("failed to fetch 30d analytics: %v", err)
	}
	stats90, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period90d)
	if err != nil {
		log.Fatalf("failed to fetch 90d analytics: %v", err)
	}
	stats365, _, err := client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period365d)
	if err != nil {
		log.Fatalf("failed to fetch 365d analytics: %v", err)
	}

	// Build lookup maps: formula name → total installs for 90d and 365d.
	count90 := make(map[string]float64)
	for _, item := range stats90.Items {
		count90[item.Formula] = parseCount(item.Count)
	}
	count365 := make(map[string]float64)
	for _, item := range stats365.Items {
		count365[item.Formula] = parseCount(item.Count)
	}

	// Take the top 50 from 30d and compute momentum.
	top := stats30.Items
	if len(top) > 50 {
		top = top[:50]
	}

	rows := make([]formulaGrowth, 0, len(top))
	for _, item := range top {
		c30 := parseCount(item.Count)
		c90 := count90[item.Formula]
		c365 := count365[item.Formula]

		rate90 := c90 / 3.0   // avg monthly over 90d
		rate365 := c365 / 12.0 // avg monthly over 365d

		momentum := 0.0
		if rate90 > 0 {
			momentum = c30 / rate90
		}

		rows = append(rows, formulaGrowth{
			name:     item.Formula,
			count30d: c30,
			rate90d:  rate90,
			rate365d: rate365,
			momentum: momentum,
		})
	}

	// Sort by momentum descending.
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].momentum > rows[j].momentum
	})

	fmt.Printf("Trending Homebrew Formulae — momentum = 30d installs ÷ avg monthly (90d)\n")
	fmt.Printf("Period: %s → %s\n\n", stats30.StartDate, stats30.EndDate)
	fmt.Printf("%-30s %12s %12s %12s %10s\n", "Formula", "30d", "90d/mo avg", "365d/mo avg", "Momentum")
	fmt.Printf("%s\n", strings.Repeat("─", 80))

	for _, r := range rows {
		trend := ""
		if r.momentum >= 1.2 {
			trend = " ↑↑"
		} else if r.momentum >= 1.05 {
			trend = " ↑"
		} else if r.momentum <= 0.8 {
			trend = " ↓↓"
		} else if r.momentum <= 0.95 {
			trend = " ↓"
		}
		fmt.Printf("%-30s %12.0f %12.0f %12.0f %9.2fx%s\n",
			r.name, r.count30d, r.rate90d, r.rate365d, r.momentum, trend)
	}
}

// parseCount converts a comma-formatted count string (e.g. "510,810") to float64.
func parseCount(s string) float64 {
	clean := strings.ReplaceAll(s, ",", "")
	v, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return v
}
