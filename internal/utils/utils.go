package utils

import (
	"github.com/joho/godotenv"
	"os"
	"solanaindexer/internal/logger"
	"strings"
)

type Config struct {
	DbName     string
	DbUsername string
	DbPassword string
	Grpc       string
}

func LoadEnvVariables() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error while getting .env file %v", err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	dbUsername := os.Getenv("MONGO_USERNAME")
	dbPassword := os.Getenv("MONGO_PASSWORD")
	grpc := os.Getenv("GRPC")

	if dbName == "" || dbUsername == "" || dbPassword == "" || grpc == "" {
		logger.Fatalf("one or more required environment variables are missing or empty")
	}

	return &Config{
		DbName:     dbName,
		DbUsername: dbUsername,
		DbPassword: dbPassword,
		Grpc:       grpc,
	}
}

func ContainsSubstring(slice []string, substring string) bool {
	for _, str := range slice {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}
