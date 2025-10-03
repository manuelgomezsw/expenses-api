package main

import (
	"expenses-api/internal/infrastructure/config"
	"expenses-api/internal/infrastructure/router"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Starting Expenses API in %s mode on port %s", cfg.Env, cfg.Port)

	// Start the application
	router.StartApp()
}
