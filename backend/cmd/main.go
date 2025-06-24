package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"r.a.w/backend/internal/api"
	"r.a.w/backend/internal/handlers"
	"r.a.w/backend/internal/router"
	"r.a.w/backend/pkg/logger"
)

func main() {
	// Initialize custom logger
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	logFilePath := filepath.Join(currentDir, "..", "logs", "backend_errors.log")
	appLogger, err := logger.NewLogger(logFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()

	// Load environment variables from .env file
	envFilePath := filepath.Join(currentDir, "..", "..", ".env")
	err = godotenv.Load(envFilePath)
	if err != nil {
		appLogger.Error("Error loading .env file: %v", err)
		return
	}

	// Get API keys from environment
	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	omdbAPIKey := os.Getenv("OMDB_API_KEY")

	if tmdbAPIKey == "" || omdbAPIKey == "" {
		appLogger.Error("TMDB_API_KEY and OMDB_API_KEY environment variables must be set")
		return
	}

	// Initialize services
	movieService := api.NewMovieService(tmdbAPIKey, omdbAPIKey, appLogger)

	// Initialize handlers
	movieHandler := handlers.NewMovieHandler(movieService, appLogger)

	// Setup routes
	r := router.SetupRoutes(movieHandler)

	// Start server
	fmt.Println("Server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
