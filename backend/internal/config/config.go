package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv                    string
	Port                      string
	DatabaseURL               string
	JWTSecret                 string
	AccessTokenExpiresMinutes int
}

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return AppConfig{
		AppEnv:                    getEnv("APP_ENV", "development"),
		Port:                      getEnv("PORT", "8080"),
		DatabaseURL:               getRequiredEnv("DATABASE_URL"),
		JWTSecret:                 getRequiredEnv("JWT_SECRET"),
		AccessTokenExpiresMinutes: getEnvAsInt("ACCESS_TOKEN_EXPIRES_MINUTES", 60),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Environment variable %s must be an integer, got: %s", key, valueStr)
	}
	return value
}
