// internal/config/config.go
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr     string
	MongoURI      string
	RateLimit     float64
	CacheLifetime time.Duration
	Environment   string
}

func Load() *Config {
	rateLimit, _ := strconv.ParseFloat(getEnvOr("RATE_LIMIT", "1.0"), 64)
	cacheMinutes, _ := strconv.Atoi(getEnvOr("CACHE_MINUTES", "5"))

	return &Config{
		ServerAddr:    getEnvOr("SERVER_ADDR", ":8080"),
		MongoURI:      getEnvOr("MONGO_URI", "mongodb://localhost:27017"),
		RateLimit:     rateLimit,
		CacheLifetime: time.Duration(cacheMinutes) * time.Minute,
		Environment:   getEnvOr("ENV", "development"),
	}
}

func getEnvOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
