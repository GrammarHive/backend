// middleware/handler.go
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grammarhive-backend/core/auth"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/grammar"
	"net/http"
	"os"
)

type APIHandler struct {
	grammarService *grammar.Service
	db            *database.MongoDB
	auth          *auth.Authenticator
}

func New(db *database.MongoDB) *APIHandler {
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	
	return &APIHandler{
		grammarService: grammar.NewService(),
		db:            db,
		auth:          auth.NewAuth0(auth0Domain, auth0Audience),
	}
}

func (h *APIHandler) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (h *APIHandler) handleOptions(w http.ResponseWriter, _ *http.Request) {
	h.setCORSHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	grammarType := r.URL.Query().Get("type")
	userID := r.URL.Query().Get("userId")

	// Validate required parameters
	if grammarType == "" || userID == "" {
		http.Error(w, "Missing required parameters: type and userId", http.StatusBadRequest)
		return
	}

	grammarContent, err := h.db.GetGrammar(r.Context(), grammarType, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generatedText, err := h.grammarService.Generate(grammarContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": generatedText,
		"status":  "OK",
	})
}

func (h *APIHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var loginRequest struct {
        Code  string `json:"code"`
        State string `json:"state"`
    }

    if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Exchange authorization code for tokens
    tokenEndpoint := fmt.Sprintf("https://%s/oauth/token", h.auth.Domain)
    payload := map[string]string{
        "grant_type":    "authorization_code",
        "client_id":     os.Getenv("AUTH0_CLIENT_ID"),
        "client_secret": os.Getenv("AUTH0_CLIENT_SECRET"),
        "code":          loginRequest.Code,
        "redirect_uri":  os.Getenv("AUTH0_CALLBACK_URL"),
    }

    payloadBytes, _ := json.Marshal(payload)
    resp, err := http.Post(tokenEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var tokenResponse struct {
        AccessToken string `json:"access_token"`
        IdToken    string `json:"id_token"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tokenResponse)
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		h.handleOptions(w, r)
		return
	}

	switch r.URL.Path {
	case "/":
		h.handleRoot(w, r)
	case "/api/login":
    	h.handleLogin(w, r)
	case "/api/generate":
		h.auth.Middleware(h.handleGenerate)(w, r)
	default:
		http.NotFound(w, r)
	}
}
