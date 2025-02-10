// core/config/config.go
package config

import (
	"os"
)

type Config struct {
    MongoURI         string
    ServerAddr       string
    Auth0Domain      string
    Auth0ClientID    string
    Auth0ClientSecret string
}

func Load() Config {
    return Config{
        MongoURI:         os.Getenv("MONGO_URI"),
        ServerAddr:       os.Getenv("SERVER_ADDR"),
        Auth0Domain:      os.Getenv("AUTH0_DOMAIN"),
        Auth0ClientID:    os.Getenv("AUTH0_CLIENT_ID"),
        Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
    }
}
