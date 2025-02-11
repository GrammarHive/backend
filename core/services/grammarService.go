// api/services/grammarService.go
package services

import (
	"context"
	"grammarhive-backend/core/database"
	"grammarhive-backend/core/grammar"
	"log"
)

type GrammarGenService struct {
	DB    *database.MongoDB
	GrammarService  *grammar.Service
}

func NewGrammarService(db *database.MongoDB) *GrammarGenService {
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	return &GrammarGenService{
		DB:              db,
		GrammarService:  grammar.NewGrammarGenService(),
	}
}

// Generate handles the logic for generating text from the grammar
func (s *GrammarGenService) Generate(grammarID string) (string, error) {
	grammarContent, err := s.DB.GetGrammar(context.Background(), grammarID)
	if err != nil {
		return "", err
	}
	return s.GrammarService.ExecuteGrammarGen(grammarContent)
}

// GenerateMultiple handles generating multiple texts
func (s *GrammarGenService) GenerateMultiple(grammarID string, count int) ([]string, error) {
	grammarContent, err := s.DB.GetGrammar(context.Background(), grammarID)
	if err != nil {
		return nil, err
	}
	return s.GrammarService.GenerateMultiple(grammarContent, count)
}
