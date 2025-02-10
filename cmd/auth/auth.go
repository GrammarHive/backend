package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load environment variables for Auth0 configuration
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	domain := os.Getenv("AUTH0_DOMAIN")
	audience := os.Getenv("AUTH0_AUDIENCE")

	if clientID == "" || clientSecret == "" || domain == "" || audience == "" {
		log.Fatal("Please set AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, AUTH0_DOMAIN, and AUTH0_AUDIENCE environment variables")
	}

	// Prepare the token request payload
	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}

	// Make the token request
	tokenURL := fmt.Sprintf("https://%s/oauth/token", domain)
	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Failed to request token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Pretty print the response
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format JSON: %v", err)
	}

	fmt.Println("Auth0 Token Response:")
	fmt.Println(prettyJSON.String())

	// Parse the response to extract just the access token
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatalf("Failed to parse token response: %v", err)
	}

	fmt.Println("\nTest curl command:")
	fmt.Printf("curl -H 'Authorization: Bearer %s' http://localhost:8080/api/grammar/generate?grammarId=YOUR_GRAMMAR_ID\n", 
		tokenResponse.AccessToken)
}
