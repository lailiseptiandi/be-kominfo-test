package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HOST string
	PORT string

	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string

	TimeZone string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		HOST:        getEnv("HOST", "localhost"),
		PORT:        getEnv("PORT", "9000"),
		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", ""),
		DB_USER:     getEnv("DB_USER", "macbookpro"),
		DB_NAME:     getEnv("DB_NAME", ""),
		DB_PASSWORD: getEnv("DB_PASSWORD", "ServBay.dev"),
		TimeZone:    getEnv("DB_TIMEZONE", "Asia/Jakarta"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
