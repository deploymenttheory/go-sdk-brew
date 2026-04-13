package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
)

func main() {
	client, err := brew.NewClient("")
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	casks, _, err := client.Casks.List(context.Background())
	if err != nil {
		log.Fatalf("failed to list casks: %v", err)
	}

	fmt.Printf("Total casks: %d\n\n", len(casks))

	// Print the first 3 as a sample.
	sample := casks
	if len(sample) > 3 {
		sample = sample[:3]
	}

	out, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}
	fmt.Println(string(out))
}
