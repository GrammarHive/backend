// core/auth/auth0.go
package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type Authenticator struct {
	Domain    string
	Audience  string
	jwks     *keyfunc.JWKS
}

func NewAuth0(domain, audience string) (*Authenticator, error) {

	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

	options := keyfunc.Options{
		RefreshInterval: time.Hour * 24,
		RefreshTimeout: time.Second * 10,
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWKS: %w", err)
	}

	return &Authenticator{
		Domain:    domain,
		Audience:  audience,
		jwks:      jwks,
	}, nil
}

func (auth *Authenticator) Middleware(next http.HandlerFunc) http.HandlerFunc {
	issuer := fmt.Sprintf("https://%s/", auth.Domain)
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "missing authorization token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, auth.jwks.Keyfunc, jwt.WithValidMethods([]string{"RS256"}))
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		if !claims.VerifyAudience(auth.Audience, true) {
			http.Error(w, "invalid audience", http.StatusUnauthorized)
			return
		}


		if !claims.VerifyIssuer(issuer, true) {
			http.Error(w, "invalid issuer", http.StatusUnauthorized)
			return
		}

		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		if !claims.VerifyNotBefore(time.Now().Unix(), true) {
			http.Error(w, "token not yet valid", http.StatusUnauthorized)
			return
		}
	
		next(w, r)
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
