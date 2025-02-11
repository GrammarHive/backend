package main

import (
	handler "grammarhive-backend/api"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := handler.NewApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      http.HandlerFunc(handler.Handler(app)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
