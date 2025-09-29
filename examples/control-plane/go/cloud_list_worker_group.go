package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	criblcontrolplane "github.com/criblio/cribl-control-plane-sdk-go"
	"github.com/criblio/cribl-control-plane-sdk-go/models/components"
	"github.com/joho/godotenv"
)

// listWorkerGroups demonstrates listing worker groups from Cribl Cloud
func listWorkerGroups() error {
	fmt.Println("Listing Cribl Worker Groups")
	fmt.Println(strings.Repeat("-", 40))

	// Load environment variables - try project root first, then current directory
	_ = godotenv.Load("../../../.env")
	_ = godotenv.Load(".env")

	// Get credentials from environment with placeholders
	orgID := getEnvOrDefault("CRIBL_ORG_ID", "your-org-id")
	clientID := getEnvOrDefault("CRIBL_CLIENT_ID", "your-client-id")
	clientSecret := getEnvOrDefault("CRIBL_CLIENT_SECRET", "your-client-secret")
	workspace := getEnvOrDefault("CRIBL_WORKSPACE_NAME", "main")

	// Check if credentials are properly set
	if orgID == "" || clientID == "" || clientSecret == "" ||
		strings.HasPrefix(orgID, "your-") ||
		strings.HasPrefix(clientID, "your-") ||
		strings.HasPrefix(clientSecret, "your-") {
		fmt.Println("Missing credentials! Set these environment variables:")
		fmt.Println("   CRIBL_ORG_ID")
		fmt.Println("   CRIBL_CLIENT_ID")
		fmt.Println("   CRIBL_CLIENT_SECRET")
		fmt.Println("\nCopy .env.example to .env and fill in your values")
		return nil
	}

	// Setup authentication
	security := components.Security{
		ClientOauth: &components.SchemeClientOauth{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     "https://login.cribl.cloud/oauth/token",
			Audience:     "https://api.cribl.cloud",
		},
	}

	// Create client
	baseURL := fmt.Sprintf("https://%s-%s.cribl.cloud/api/v1", workspace, orgID)
	client := criblcontrolplane.New(baseURL, criblcontrolplane.WithSecurity(security))

	// List worker groups
	fmt.Println("Fetching worker groups...")

	ctx := context.Background()

	response, err := client.Groups.List(ctx, components.ProductsCoreStream, nil)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	// Handle the response
	if response.Object != nil && response.Object.Items != nil {
		items := response.Object.Items
		fmt.Printf("\nFound %d worker group(s):\n", len(items))
		fmt.Println()

		for _, group := range items {
			groupID := group.ID
			if groupID == "" {
				groupID = "Unknown"
			}

			fmt.Printf("Worker Group: %s\n", groupID)
			fmt.Println(strings.Repeat("-", len(groupID)+16))

			// Print available fields (this is a simplified version)
			if group.ID != "" {
				fmt.Printf("   id: %s\n", group.ID)
			}
			if group.Description != nil && *group.Description != "" {
				fmt.Printf("   description: %s\n", *group.Description)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No worker groups found")
	}

	return nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	if err := listWorkerGroups(); err != nil {
		log.Fatal(err)
	}
}
