package services

import (
	"context"
	"grammarhive-backend/core/database"
	"log"
)

type ProfileService struct {
	DB    *database.MongoDB
}

func NewProfileService(db *database.MongoDB) *ProfileService {
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	return &ProfileService{
		DB: db,
	}
}

func (p *ProfileService) UploadGrammarToProfile(ctx context.Context, input *database.Grammar) error {
	return p.DB.StoreGrammar(ctx, input.GrammarID, input.Name, input.Username, input.Content, input.Version)
}
