// core/auth/auth0.go
package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

type Authenticator struct {
	Domain    string
	Audience  string
}

func NewAuth0(domain, audience string) *Authenticator {
	return &Authenticator{
		Domain:    domain,
		Audience:  audience,
	}
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func (a *Authenticator) Middleware(next http.HandlerFunc) http.HandlerFunc {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := a.Audience
			// Checks if the token was intended for API
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}

			// Verify 'iss' claim
			iss := "https://" + a.Domain + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			// Gets Auth0's public key (certificate) to verify the token's signature
			cert, err := a.getPemCert(token)
			if err != nil {
				return nil, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		if err := jwtMiddleware.CheckJWT(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (a *Authenticator) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	// Fetches JWKS (JSON Web Key Set) from Auth0
	resp, err := http.Get("https://" + a.Domain + "/.well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kty string   `json:"kty"`
			Kid string   `json:"kid"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}

	return cert, nil
}
