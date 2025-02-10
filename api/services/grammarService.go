// api/services/grammarService.go
package services

import (
	"context"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/grammar"
)

// Service holds the dependencies required for the grammar logic
type Service struct {
	DB    *database.MongoDB
	GrammarService  *grammar.Service
}

// NewGrammarService creates a new instance of Service
func NewGrammarService(db *database.MongoDB) *Service {
	return &Service{
		DB:              db,
		GrammarService:  grammar.NewService(),
	}
}

// Generate handles the logic for generating text from the grammar
func (s *Service) Generate(grammarID string) (string, error) {
	grammarContent, err := s.DB.GetGrammar(context.Background(), grammarID)
	if err != nil {
		return "", err
	}
	return s.GrammarService.Generate(grammarContent)
}

// GenerateMultiple handles generating multiple texts
func (s *Service) GenerateMultiple(grammarID string, count int) ([]string, error) {
	grammarContent, err := s.DB.GetGrammar(context.Background(), grammarID)
	if err != nil {
		return nil, err
	}
	return s.GrammarService.GenerateMultiple(grammarContent, count)
}
