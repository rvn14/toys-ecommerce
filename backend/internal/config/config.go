package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ApConfig struct {
	AppEnv string
	Port   string
}

func LoadConfig() ApConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return ApConfig{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
