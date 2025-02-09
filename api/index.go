// api/index.go
package handler

import (
	"context"
	"net/http"
	"time"

	api "grammarhive-backend/components"
	config "grammarhive-backend/components/config"
	database "grammarhive-backend/components/database"
)
func Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cfg := config.Load()
	mongoClient, err := database.NewMongoDB(ctx, cfg.MongoURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoClient.Close(ctx)

	h := api.New(mongoClient)
	h.ServeHTTP(w, r)
}
