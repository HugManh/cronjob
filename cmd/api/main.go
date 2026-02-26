package main

import (
	"github.com/joho/godotenv"

	"github.com/HugManh/cronjob/internal/startup"
	"github.com/HugManh/cronjob/pkg/logger"
)

func init() {
	// Initialize logger
	godotenv.Load()
	logger.NewLogger(&logger.Config{})
}

func main() {
	startup.Server()
}
