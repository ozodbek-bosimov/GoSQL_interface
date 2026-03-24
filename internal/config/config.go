package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DatabaseURL   string
	ServerPort    string
	SessionSecret string
	Environment   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gosql_db?sslmode=disable"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		SessionSecret: getEnv("SESSION_SECRET", "change-this-secret-key"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}

	return config
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
