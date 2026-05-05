package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerPort    string
	SessionSecret string
	Environment   string
}

func LoadConfig() *Config {
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
