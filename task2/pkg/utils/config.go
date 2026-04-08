package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvCollectorPort = "COLLECTOR_PORT"
	EnvCollectorAddr = "COLLECTOR_ADDR"
	EnvHTTPPort      = "HTTP_PORT"
)

const (
	DefaultCollectorPort = ":50051"
	DefaultCollectorAddr = "localhost:50051"
	DefaultHTTPPort      = ":8080"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: could not load .env file: %v. Using system environment variables.", err)
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
