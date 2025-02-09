// components/config/config.go
package config

import (
	"os"
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

	return &Config{
		ServerAddr:    os.Getenv("SERVER_ADDR"),
		MongoURI:      os.Getenv("MONGO_URI"),
		Environment:   os.Getenv("ENV"),
	}
}
