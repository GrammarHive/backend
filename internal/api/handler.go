// internal/handler/handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.resumes.guide/internal/database"
	"go.resumes.guide/internal/grammar"
)

type Handler struct {
	grammarService *grammar.Service
	db            *database.MongoDB
}

func New(db *database.MongoDB) *Handler {
	return &Handler{
		grammarService: grammar.NewService(),
		db:            db,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	grammarContent, err := h.db.GetGrammar(r.Context(), "resume")
	if err != nil {
		http.Error(w, "Failed to fetch grammar", http.StatusInternalServerError)
		return
	}

	generatedText, err := h.grammarService.Generate(grammarContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": generatedText,
		"status":  "OK",
	})
}
