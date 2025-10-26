package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	criblcloudmanagementsdkgo "github.com/criblio/cribl-cloud-management-sdk-go"
	"github.com/criblio/cribl-cloud-management-sdk-go/models/components"
	"github.com/joho/godotenv"
)

// listWorkspaces demonstrates listing workspaces from Cribl Management Plane
func listWorkspaces() error {
	fmt.Println("Listing Cribl Workspaces")
	fmt.Println(strings.Repeat("-", 40))

	// Load environment variables - try project root first, then current directory
	_ = godotenv.Load("../../../.env")
	_ = godotenv.Load(".env")

	// Get credentials from environment with placeholders
	orgID := getEnvOrDefault("CRIBL_ORG_ID", "your-org-id")
	clientID := getEnvOrDefault("CRIBL_CLIENT_ID", "your-client-id")
	clientSecret := getEnvOrDefault("CRIBL_CLIENT_SECRET", "your-client-secret")
	domain := getEnvOrDefault("CRIBL_DOMAIN", "cribl.cloud")

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

	// Create client with authentication
	client := criblcloudmanagementsdkgo.New(
		criblcloudmanagementsdkgo.WithServerURL(fmt.Sprintf("https://gateway.%s", domain)),
		criblcloudmanagementsdkgo.WithSecurity(components.Security{
			ClientOauth: &components.SchemeClientOauth{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				TokenURL:     fmt.Sprintf("https://login.%s/oauth/token", domain),
				Audience:     fmt.Sprintf("https://api.%s", domain),
			},
		}),
	)

	// List workspaces
	fmt.Println("Fetching workspaces...")

	ctx := context.Background()
	response, err := client.Workspaces.List(ctx, orgID)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	// Handle the response
	if response.WorkspacesListResponseDTO != nil {
		items := response.WorkspacesListResponseDTO.Items
		fmt.Printf("\nFound %d workspace(s):\n", len(items))
		fmt.Println()

		for _, workspace := range items {
			workspaceID := workspace.WorkspaceID
			if workspaceID == "" {
				workspaceID = "Unknown"
			}

			fmt.Printf("Workspace: %s\n", workspaceID)
			fmt.Println(strings.Repeat("-", len(workspaceID)+12))

			// Print available fields
			if workspace.WorkspaceID != "" {
				fmt.Printf("   workspaceId: %s\n", workspace.WorkspaceID)
			}
			if workspace.Region != "" {
				fmt.Printf("   region: %s\n", string(workspace.Region))
			}
			if !workspace.LastUpdated.IsZero() {
				fmt.Printf("   lastUpdated: %s\n", workspace.LastUpdated.String())
			}
			if workspace.LeaderFQDN != "" {
				fmt.Printf("   leaderFQDN: %s\n", workspace.LeaderFQDN)
			}
			if workspace.State != "" {
				fmt.Printf("   state: %s\n", string(workspace.State))
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No workspaces found")
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
	if err := listWorkspaces(); err != nil {
		log.Fatal(err)
	}
}
