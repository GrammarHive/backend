// api/index.go
package handler

import (
	"context"
	"net/http"
	"time"

	api "go.resumes.guide/internal/api"
	config "go.resumes.guide/internal/config"
	database "go.resumes.guide/internal/database"
)
func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cfg := config.Load()
	mongoClient, err := database.NewMongoDB(ctx, cfg.MongoURI)
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoClient.Close(ctx)

	h := api.New(mongoClient)
	h.ServeHTTP(w, r)
}
