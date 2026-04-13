// Recipe: Cask Catalog Summary
//
// Loads the full cask catalog and produces a summary report broken down by:
//   - Total cask count
//   - Casks with auto-update enabled vs manual
//   - Casks that are deprecated or disabled
//   - Top taps by cask count
//   - Casks with Rosetta caveats (Apple Silicon compatibility notes)
//
// Then cross-references with 30-day cask install analytics to annotate
// the most-downloaded casks with their catalog metadata.
package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/casks"
)

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Load the full cask catalog.
	catalog, _, err := client.Casks.List(ctx)
	if err != nil {
		log.Fatalf("failed to list casks: %v", err)
	}

	// Step 2: Compute catalog statistics.
	var (
		autoUpdate  int
		deprecated  int
		disabled    int
		rosetta     int
		tapCounts   = make(map[string]int)
	)

	for _, c := range catalog {
		if c.AutoUpdates {
			autoUpdate++
		}
		if c.Deprecated {
			deprecated++
		}
		if c.Disabled {
			disabled++
		}
		if c.CaveatsRosetta != nil && *c.CaveatsRosetta {
			rosetta++
		}
		tapCounts[c.Tap]++
	}

	// Sort taps by count.
	type tapEntry struct {
		tap   string
		count int
	}
	tapList := make([]tapEntry, 0, len(tapCounts))
	for tap, count := range tapCounts {
		tapList = append(tapList, tapEntry{tap, count})
	}
	sort.Slice(tapList, func(i, j int) bool {
		return tapList[i].count > tapList[j].count
	})

	fmt.Printf("╔══════════════════════════════════════════════════════╗\n")
	fmt.Printf("║           Homebrew Cask Catalog Summary              ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("Total casks:        %d\n", len(catalog))
	fmt.Printf("Auto-update:        %d (%.1f%%)\n", autoUpdate, pct(autoUpdate, len(catalog)))
	fmt.Printf("Deprecated:         %d\n", deprecated)
	fmt.Printf("Disabled:           %d\n", disabled)
	fmt.Printf("Rosetta caveats:    %d\n", rosetta)

	fmt.Printf("\nTop taps by cask count:\n")
	for i, t := range tapList {
		if i >= 5 {
			break
		}
		fmt.Printf("  %-35s %d casks\n", t.tap, t.count)
	}

	// Step 3: Fetch 30d cask install analytics and cross-reference.
	installStats, _, err := client.Analytics.ListCaskInstalls(ctx, analytics.Period30d)
	if err != nil {
		log.Fatalf("failed to fetch cask analytics: %v", err)
	}

	// Build a token → cask metadata map for quick lookup.
	caskByToken := make(map[string]casks.ResourceCask, len(catalog))
	for _, c := range catalog {
		caskByToken[c.Token] = c
	}

	// Collect entries with install counts, sort by count descending.
	type installEntry struct {
		token string
		count string
		cask  *casks.ResourceCask
	}
	entries := make([]installEntry, 0, len(installStats.Formulae))
	for token, items := range installStats.Formulae {
		if len(items) == 0 {
			continue
		}
		entry := installEntry{token: token, count: items[0].Count}
		if c, ok := caskByToken[token]; ok {
			entry.cask = &c
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return parseCommaInt(entries[i].count) > parseCommaInt(entries[j].count)
	})

	fmt.Printf("\nTop 10 most-installed casks (30d) with catalog metadata:\n")
	fmt.Printf("%s\n", strings.Repeat("─", 80))
	fmt.Printf("%-25s %10s  %-10s  %s\n", "Token", "Installs", "Version", "Description")
	fmt.Printf("%s\n", strings.Repeat("─", 80))

	for i, e := range entries {
		if i >= 10 {
			break
		}
		version := "—"
		desc := "—"
		if e.cask != nil {
			version = e.cask.Version
			desc = truncate(e.cask.Desc, 35)
		}
		fmt.Printf("%-25s %10s  %-10s  %s\n", e.token, e.count, version, desc)
	}
}

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func parseCommaInt(s string) int {
	clean := strings.ReplaceAll(s, ",", "")
	n := 0
	fmt.Sscanf(clean, "%d", &n)
	return n
}
