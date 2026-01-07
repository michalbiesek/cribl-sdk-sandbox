package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	criblcontrolplane "github.com/criblio/cribl-control-plane-sdk-go"
	"github.com/criblio/cribl-control-plane-sdk-go/models/components"
	"github.com/joho/godotenv"
)

// captureEvents demonstrates calling system.captures.get on on-premise Cribl Control Plane
func captureEvents() error {
	fmt.Println("Calling system.captures.get on On-Premise Cribl")
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

	// Call system.captures.get
	fmt.Println("\nCalling system.captures.get...")
	fmt.Println("Payload: {\"filter\":\"__inputId=='datagen:datagenTest'\",\"duration\":10,\"maxEvents\":10,\"level\":0}")

	captureParams := components.CaptureParams{
		Filter:    "__inputId=='datagen:datagenTest'",
		Duration:  10,
		MaxEvents: 10,
		Level:     components.CaptureLevelZero,
	}

	capturesResponse, err := client.System.Captures.Get(ctx, captureParams)
	if err != nil {
		return fmt.Errorf("Error calling system.captures.get: %w", err)
	}

	fmt.Println("\nCaptures response:")
	fmt.Println(strings.Repeat("-", 50))
	if capturesResponse != nil && capturesResponse.HTTPMeta.Response != nil {
		// Read the stream response (JSONL format)
		body := capturesResponse.HTTPMeta.Response.Body
		if body != nil {
			defer body.Close()

			scanner := bufio.NewScanner(body)
			eventCount := 0

			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}

				eventCount++
				fmt.Printf("Event %d:\n", eventCount)

				// Try to parse and pretty-print JSON
				var eventData map[string]interface{}
				if err := json.Unmarshal([]byte(line), &eventData); err == nil {
					jsonBytes, _ := json.MarshalIndent(eventData, "  ", "  ")
					fmt.Println(string(jsonBytes))
				} else {
					// If not valid JSON, just print the line
					fmt.Println(line)
				}
				fmt.Println()
			}

			if err := scanner.Err(); err != nil && err != io.EOF {
				return fmt.Errorf("Error reading captures stream: %w", err)
			}

			if eventCount == 0 {
				fmt.Println("No events captured")
			}
		} else {
			fmt.Println("Response body is nil")
		}
	} else {
		fmt.Println("No response received or invalid response structure")
		if capturesResponse != nil {
			fmt.Printf("Response structure: %+v\n", capturesResponse)
		}
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
	if err := captureEvents(); err != nil {
		log.Fatal(err)
	}
}
