package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/services"
	"io"
	"net/http"
	"time"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(dbService *database.MongoDB) *ProfileHandler {
	return &ProfileHandler{
		profileService: services.NewProfileService(dbService),
	}
}

// GenerateRandomID generates a random hex string of specified length
func GenerateRandomID(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func (p *ProfileHandler) HandleGetGrammarByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	err := p.profileService.ValidateUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grammars, err := p.profileService.DB.GetGrammarsByUsername(username)
	if err != nil {
		http.Error(w, "Error retrieving grammar entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grammars)
}

func (p *ProfileHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	username := r.FormValue("username")

	grammarID, err := GenerateRandomID(6)
	if err != nil {
		http.Error(w, "Error generating random value", http.StatusInternalServerError)
	}

	file, _, err := r.FormFile("grammarFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validate all user generated content here
	err = p.profileService.ValidateInput(name, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate user file here
	if err := p.profileService.ValidateFile(file); err != nil {
		http.Error(w, fmt.Sprintf("File validation failed: %v", err), http.StatusBadRequest)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	currentTime := time.Now()

	input := &database.Grammar{
		GrammarID: grammarID,
		Name:      name,
		Version:   0,
		Content:   string(content),
		Username:  username,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	if err = p.profileService.UploadGrammarToProfile(context.Background(), input); err != nil {
		http.Error(w, "Error storing grammar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "File uploaded and grammar stored successfully!",
		"status":  "success",
	})
}
