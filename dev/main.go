// dev/main.go
package main

import (
	"log"
	"net/http"

	"go.resumes.guide/api"
)

func main() {
	// Mount the serverless function handler
	http.HandleFunc("/api/generate", api.Handler)
	
	// Run the server
	port := ":4000"
	log.Printf("Starting server on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
