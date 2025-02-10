// api/auth/login.go
package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenEndpoint := fmt.Sprintf("https://%s/oauth/token", os.Getenv("AUTH0_DOMAIN"))
	payload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     os.Getenv("AUTH0_CLIENT_ID"),
		"client_secret": os.Getenv("AUTH0_CLIENT_SECRET"),
		"audience":      os.Getenv("AUTH0_AUDIENCE"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to create token request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(tokenEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to get token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse)
}
