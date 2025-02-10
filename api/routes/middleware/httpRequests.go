// middleware/middleware.go
package handler

import (
	"net/http"
)

// SetCORSHeaders sets the CORS headers for the HTTP response.
func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// HandleOptions handles OPTIONS request method
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	SetCORSHeaders(w)
	w.WriteHeader(http.StatusOK)
}
