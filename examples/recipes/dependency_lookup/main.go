// Recipe: Dependency Lookup
//
// Given a starting formula, recursively fetches and prints its full dependency
// tree — both runtime and build dependencies — up to a configurable depth.
// Useful for understanding the transitive dependency closure of a formula.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
)

const (
	startFormula = "ffmpeg" // Replace with any formula name.
	maxDepth     = 3
)

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	visited := make(map[string]bool)

	fmt.Printf("Dependency tree for %q (max depth: %d)\n\n", startFormula, maxDepth)
	walk(ctx, client, startFormula, 0, visited)
}

func walk(ctx context.Context, client *brew.Client, name string, depth int, visited map[string]bool) {
	if depth > maxDepth || visited[name] {
		return
	}
	visited[name] = true

	prefix := strings.Repeat("  ", depth)

	formula, _, err := client.Formulae.GetByName(ctx, name)
	if err != nil {
		fmt.Printf("%s├─ %s [error: %v]\n", prefix, name, err)
		return
	}

	marker := "├─"
	if depth == 0 {
		marker = "◉"
	}
	fmt.Printf("%s%s %s (%s)  %s\n", prefix, marker, formula.Name, formula.Versions.Stable, formula.Desc)

	// Print build dependencies as annotations (not recursed into, to avoid noise).
	if len(formula.BuildDependencies) > 0 {
		fmt.Printf("%s   build-deps: %s\n", prefix, joinStrings(formula.BuildDependencies))
	}

	for _, dep := range formula.Dependencies {
		walk(ctx, client, dep, depth+1, visited)
	}
}

func joinStrings(ss []string) string {
	return strings.Join(ss, ", ")
}
