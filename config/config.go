package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseUrl string
	ServerPort  string
	GinMode     string
	JwtSecret   string
}

func Load() *Config {
	cfg := &Config{
		DatabaseUrl: getEnv("DATABASE_URL", ""),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "release"),
		JwtSecret:   getEnv("JWT_SECRET", "defaultsecret"),
	}
	if v, _ := os.LookupEnv("SILENCE_LOGS"); v != "true" {
		log.Printf("Configuration loaded: %+v\n", cfg)
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultValue
}
