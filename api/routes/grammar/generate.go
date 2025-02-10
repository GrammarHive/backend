package handler

import (
	"encoding/json"
	"fmt"
	"grammarhive-backend/api/routes/validation"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/services"
	"net/http"
)

var grammarService *services.Service

func Init(dbService *database.MongoDB) {
	grammarService = services.NewGrammarService(dbService)
}

// HandleGenerate handles the generation of a single grammar text
func HandleGenerate(w http.ResponseWriter, r *http.Request) {
	grammarID, err := validation.ValidateGenerateRequest(r)
	if err != nil {
		http.Error(w, "Missing required parameter: grammarId", http.StatusBadRequest)
		return
	}

	// Utilize the GrammarService to generate the text
	generatedText, err := grammarService.Generate(grammarID)
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
func HandleGenerateList(w http.ResponseWriter, r *http.Request) {
	grammarID, count, err := validation.ValidateGenerateListRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Use GrammarService to generate multiple texts
	messages, err := grammarService.GenerateMultiple(grammarID, count)
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
