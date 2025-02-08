// dev/main.go
package main

import (
	"log"
	"net/http"

	handler "go.resumes.guide/api"
)

func main() {
	// Mount the serverless function handler
	http.HandleFunc("/api/generate", handler.Handler)
	
	// Run the server
	port := ":4000"
	log.Printf("Starting server on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
