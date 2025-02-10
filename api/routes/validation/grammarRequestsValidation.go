// api/validation.go
package validation

import (
	"net/http"
	"strconv"
)

// ValidateGenerateRequest validates the request for generating results
func ValidateGenerateRequest(r *http.Request) (string, error) {
	grammarID := r.URL.Query().Get("grammarId")
	if grammarID == "" {
		return "", http.ErrMissingFile
	}
	return grammarID, nil
}

// ValidateGenerateListRequest validates the request for generating multiple results
func ValidateGenerateListRequest(r *http.Request) (string, int, error) {
	grammarID := r.URL.Query().Get("grammarId")
	if grammarID == "" {
		return "", 0, http.ErrMissingFile
	}

	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		return "", 0, http.ErrNotSupported
	}
	return grammarID, count, nil
}
