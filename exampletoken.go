package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func sam() {
	// Load the OAuth2 credentials file
	credentialsFile := "credentials.json"

	// Set up OAuth2 configuration
	oauthConfig, err := getOAuthConfigs(credentialsFile)
	if err != nil {
		log.Fatalf("Failed to get OAuth2 config: %v", err)
	}

	// Obtain an access token using the authentication flow
	ctx := context.Background()
	token, err := getToken(ctx, oauthConfig)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Create a new HTTP client with the access token
	httpClient := oauthConfig.Client(ctx, token)

	// Create a new Google Sheets service
	service, err := sheets.New(httpClient)
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
	}

	// Make API calls using the service
	// ...

	fmt.Println("Google Sheets integration set up successfully")
}

func getOAuthConfigs(credentialsFile string) (*oauth2.Config, error) {
	// Load the OAuth2 credentials from file
	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	// Parse the credentials file
	config, err := google.ConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	return config, nil
}

func getTokens(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	// Perform the authentication flow to obtain an access token
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser to authenticate:\n%s\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("failed to read authorization code: %v", err)
	}

	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	return token, nil
}
