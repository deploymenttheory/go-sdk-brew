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

	name := "wget" // Replace with the desired formula name.

	formula, _, err := client.Formulae.GetByName(context.Background(), name)
	if err != nil {
		log.Fatalf("failed to get formula %q: %v", name, err)
	}

	out, err := json.MarshalIndent(formula, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal output: %v", err)
	}
	fmt.Println(string(out))
}
