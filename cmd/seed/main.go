// cmd/seed/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"grammarhive-backend/core/config"
	"grammarhive-backend/core/database"
)

func main() {
	var grammarURL string
	flag.StringVar(&grammarURL, "url", "https://raw.githubusercontent.com/HarryZ10/api.resumes.guide/main/static/resume.g", "URL of the grammar file")
	flag.Parse()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := database.NewMongoDB(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close(ctx)

	grammarContent, err := fetchGrammar(grammarURL)
	if err != nil {
		log.Fatalf("Failed to fetch grammar: %v", err)
	}

	if err := db.StoreGrammar(ctx, "resume", "admin", grammarContent); err != nil {
		log.Fatalf("Failed to store grammar: %v", err)
	}

	log.Println("Successfully seeded grammar into database")
}

func fetchGrammar(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
