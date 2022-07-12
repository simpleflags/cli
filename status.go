package main

import (
	"fmt"
	"log"
	"os"
)

type statusCommand struct {
}

func (c statusCommand) Execute(_ []string) error {
	status := "healthy"
	fmt.Println("--- Environment variables ---")
	fmt.Printf("Admin server URL: %s\n", api.BaseURL())
	fmt.Printf("Default account: %s\n", os.Getenv("SF_ACCOUNT"))
	fmt.Printf("Default project: %s\n", os.Getenv("SF_PROJECT"))
	fmt.Println("--- Server ---")
	fmt.Printf("Status: %s\n", status)
	fmt.Println("--- API KEY ---")
	fmt.Printf("API Key: %s\n", os.Getenv("SF_API_KEY"))
	fmt.Printf("API Key type: %s\n", os.Getenv("SF_KEY_TYPE"))
	return nil
}

func init() {
	sc := statusCommand{}
	_, err := parser.AddCommand(
		"status",
		"Show status of SF and environment variables",
		"Show status of SF and environment variables",
		&sc,
	)

	if err != nil {
		log.Printf("error adding command %v", err)
	}
}
