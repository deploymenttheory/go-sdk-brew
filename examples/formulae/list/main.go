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

	formulae, _, err := client.Formulae.List(context.Background())
	if err != nil {
		log.Fatalf("failed to list formulae: %v", err)
	}

	fmt.Printf("Total formulae: %d\n\n", len(formulae))

	// Print the first 3 as a sample.
	sample := formulae
	if len(sample) > 3 {
		sample = sample[:3]
	}

	out, err := json.MarshalIndent(sample, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}
	fmt.Println(string(out))
}
