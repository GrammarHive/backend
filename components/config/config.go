// components/config/config.go
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr     string
	MongoURI      string
	RateLimit     float64
	CacheLifetime time.Duration
	Environment   string
}

func Load() *Config {

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }

	return &Config{
		ServerAddr:    getEnvOr("SERVER_ADDR", ":8080"),
		MongoURI:      getEnvOr("MONGO_URI", "mongodb://localhost:27017"),
		Environment:   getEnvOr("ENV", "development"),
	}
}

func getEnvOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
