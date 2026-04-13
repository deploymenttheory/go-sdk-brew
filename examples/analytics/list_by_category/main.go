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

	// Category: CategoryInstall, CategoryInstallOnRequest, or CategoryBuildError.
	// Period:   Period30d, Period90d, or Period365d.
	category := analytics.CategoryInstall
	period := analytics.Period30d

	result, _, err := client.Analytics.ListByCategory(context.Background(), category, period)
	if err != nil {
		log.Fatalf("failed to list analytics: %v", err)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}
	fmt.Println(string(out))
}
