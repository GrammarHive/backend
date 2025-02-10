// middleware/handler.go
package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grammarhive-backend/core/auth"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/grammar"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceHandler struct {
	grammarService   *grammar.Service
	dbService        *database.MongoDB
	authService      *auth.Authenticator
}

func New(dbService *database.MongoDB) *ServiceHandler {
	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	authService, err := auth.NewAuth0(auth0Domain, auth0Audience)
	if err != nil {
		log.Fatalf("Failed to initialize authenticator: %v", err)
	}

	return &ServiceHandler{
		grammarService:   grammar.NewService(),
		dbService:        dbService,
		authService:      authService,
	}
}

func (router *ServiceHandler) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (router *ServiceHandler) handleOptions(w http.ResponseWriter, _ *http.Request) {
	router.setCORSHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (router *ServiceHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (router *ServiceHandler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	grammarID := r.URL.Query().Get("grammarId")

	// Validate required parameters
	if grammarID == "" {
		http.Error(w, "Missing required parameter: grammarId", http.StatusBadRequest)
		return
	}

	// Get grammar content from MongoDB
	grammarContent, err := router.dbService.GetGrammar(r.Context(), grammarID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Grammar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch grammar", http.StatusInternalServerError)
		return
	}

	// Generate text using the grammar
	generatedText, err := router.grammarService.Generate(grammarContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the generated text
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": generatedText,
		"status": "success",
		"grammarId": grammarID,
	})
}

func (router *ServiceHandler) handleGenerateList(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	grammarID := r.URL.Query().Get("grammarId")

	// Validate required parameters
	if grammarID == "" {
		http.Error(w, "Missing required parameter: grammarId", http.StatusBadRequest)
		return
	}

	// Get grammar content from db
	grammarContent, err := router.dbService.GetGrammar(r.Context(), grammarID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Grammar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch grammar", http.StatusInternalServerError)
		return
	}

	// Generate multiple texts using the grammar
	messages, err := router.grammarService.GenerateMultiple(grammarContent, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the generated texts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"count": len(messages),
		"status": "success",
		"grammarId": grammarID,
	})
}


func (router *ServiceHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token using client credentials
	tokenEndpoint := fmt.Sprintf("https://%s/oauth/token", router.authService.Domain)
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
		TokenType  string `json:"token_type"`
		ExpiresIn  int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Failed to parse token response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResponse)
}



func (router *ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.setCORSHeaders(w)

	if r.Method == "OPTIONS" {
		router.handleOptions(w, r)
		return
	}

	switch r.URL.Path {
	case "/":
		router.handleRoot(w, r)
	case "/api/login":
		router.handleLogin(w, r)
	// case "/api/grammar/upload":
	// 	if r.Method == http.MethodPost {
	// 		h.auth.Middleware(h.handleGrammarUpload)(w, r)
	// 	} else {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	case "/api/grammar/generate":
		if r.Method == http.MethodGet {
			if r.URL.Query().Get("list") != "" {
				router.authService.Middleware(router.handleGenerateList)(w, r)
			} else {
				router.authService.Middleware(router.handleGenerate)(w, r)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}
