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

	name := "iterm2" // Replace with the desired cask token.

	cask, _, err := client.Casks.GetByName(context.Background(), name)
	if err != nil {
		log.Fatalf("failed to get cask %q: %v", name, err)
	}

	out, err := json.MarshalIndent(cask, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}
	fmt.Println(string(out))
}
