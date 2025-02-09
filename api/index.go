// api/index.go
package api

import (
	"context"
	"net/http"
	"time"

	config "go.resumes.guide/api/config"
	database "go.resumes.guide/api/database"
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

	h := New(mongoClient)
	h.ServeHTTP(w, r)
}
