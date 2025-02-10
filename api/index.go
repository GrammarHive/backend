// api/index.go
package handler

import (
	"context"
	"net/http"
	"time"

	config "grammarhive-backend/core/config"
	database "grammarhive-backend/core/database"
	server "grammarhive-backend/middleware"
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

	h := server.New(mongoClient)
	h.ServeHTTP(w, r)
}
