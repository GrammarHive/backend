package handler

import (
	"encoding/json"
	"fmt"
	"grammarhive-backend/api/routes/validation"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/services"
	"net/http"
)

type GrammarHandler struct {
	grammarService *services.Service
}

func NewGrammarHandler(dbService *database.MongoDB) *GrammarHandler {
	return &GrammarHandler{
		grammarService: services.NewGrammarService(dbService),
	}
}

// HandleGenerate handles the generation of a single grammar text
func (h *GrammarHandler) HandleGenerate(w http.ResponseWriter, r *http.Request) {
	// Check if grammarService is initialized
	if h.grammarService == nil {
		http.Error(w, "Grammar service is not initialized", http.StatusInternalServerError)
		return
	}

	grammarID, err := validation.ValidateGenerateRequest(r)
	if err != nil {
		http.Error(w, "Missing required parameter: grammarId", http.StatusBadRequest)
		return
	}

	// Utilize the GrammarService to generate the text
	generatedText, err := h.grammarService.Generate(grammarID)
	if err != nil {
		http.Error(w, "Generation failed", http.StatusInternalServerError)
		return
	}

	// Return the generated text
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    generatedText,
		"status":     "success",
		"grammarId":  grammarID,
	})
}

// HandleGenerateList handles the generation of multiple grammar texts
func (h *GrammarHandler) HandleGenerateList(w http.ResponseWriter, r *http.Request) {
	grammarID, count, err := validation.ValidateGenerateListRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use GrammarService to generate multiple texts
	messages, err := h.grammarService.GenerateMultiple(grammarID, count)
	if err != nil {
		http.Error(w, fmt.Sprintf("Generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the generated texts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages":   messages,
		"count":      len(messages),
		"status":     "success",
		"grammarId":  grammarID,
	})
}
