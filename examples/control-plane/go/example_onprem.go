package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	criblcontrolplane "github.com/criblio/cribl-control-plane-sdk-go"
	"github.com/criblio/cribl-control-plane-sdk-go/models/components"
	"github.com/joho/godotenv"
)

// listWorkerGroups demonstrates listing worker groups from on-premise Cribl Control Plane
func listWorkerGroups() error {
	fmt.Println("Listing On-Premise Cribl Worker Groups")
	fmt.Println(strings.Repeat("-", 45))

	// Load environment variables - try project root first, then current directory
	_ = godotenv.Load("../../../.env")
	_ = godotenv.Load(".env")

	// Get credentials from environment with placeholders
	serverURL := getEnvOrDefault("CRIBL_SERVER_URL", "http://localhost:19000")
	username := getEnvOrDefault("CRIBL_USERNAME", "admin")
	password := getEnvOrDefault("CRIBL_PASSWORD", "admin")
	// Default to true for on-prem development environments with self-signed certs
	insecureTLS := getEnvOrDefault("CRIBL_INSECURE_TLS", "true") != "false"

	// Check if server URL is properly set
	if strings.HasPrefix(serverURL, "your-") {
		fmt.Println("Invalid server URL! Set this environment variable:")
		fmt.Println("   CRIBL_SERVER_URL")
		fmt.Println("\nCopy .env.example to .env and fill in your values")
		return nil
	}

	// Create base URL for API
	baseURL := strings.TrimSuffix(serverURL, "/") + "/api/v1"
	fmt.Printf("Connecting to: %s\n", baseURL)

	// Configure TLS for HTTPS with self-signed certificates
	var clientOpts []criblcontrolplane.SDKOption
	if strings.HasPrefix(serverURL, "https") && insecureTLS {
		fmt.Println("⚠️  Accepting self-signed certificates (insecure mode)")
		clientOpts = append(clientOpts, criblcontrolplane.WithClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}))
	}

	// First, create an unauthenticated client to get a token
	client := criblcontrolplane.New(baseURL, clientOpts...)

	// Authenticate with username/password to get token
	fmt.Println("Authenticating with username/password...")

	ctx := context.Background()

	loginInfo := components.LoginInfo{
		Username: username,
		Password: password,
	}

	authResponse, err := client.Auth.Tokens.Get(ctx, loginInfo)
	if err != nil {
		errMsg := err.Error()
		// Check if error is related to self-signed certificates
		if strings.Contains(errMsg, "certificate") || strings.Contains(errMsg, "x509") || strings.Contains(errMsg, "TLS") {
			fmt.Println("\n💡 Tip: If you're using a self-signed certificate, set:")
			fmt.Println("   CRIBL_INSECURE_TLS=true")
			fmt.Println("   (Only use this in development/testing environments!)")
		}
		return fmt.Errorf("Authentication failed: %w", err)
	}

	var token string
	if authResponse.AuthToken != nil {
		token = authResponse.AuthToken.Token
		fmt.Printf("Authenticated with on-prem server, token: %s\n", token)
	} else {
		return fmt.Errorf("No token received from authentication")
	}

	// Create authenticated SDK client with bearer token
	security := components.Security{
		BearerAuth: &token,
	}

	clientOpts = append(clientOpts, criblcontrolplane.WithSecurity(security))
	client = criblcontrolplane.New(baseURL, clientOpts...)
	fmt.Println("Cribl SDK client created for on-prem server")

	// List worker groups
	fmt.Println("Fetching worker groups...")

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
