package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv" // Import godotenv

	"r.a.w/backend/internal/api"
	"r.a.w/backend/pkg/logger"
)

// Serve static files from the "frontend/public" directory
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

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("../../frontend/public"))
	r.PathPrefix("/").Handler(fs)

	// Load environment variables from .env file
	envFilePath := filepath.Join(currentDir, "..", "..", ".env")
	err = godotenv.Load(envFilePath)
	if err != nil {
		appLogger.Error("Error loading .env file: %v", err)
		return // Exit if .env cannot be loaded
	}

	// Initialize API clients
	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	omdbAPIKey := os.Getenv("OMDB_API_KEY")

	if tmdbAPIKey == "" || omdbAPIKey == "" {
		appLogger.Error("TMDB_API_KEY and OMDB_API_KEY environment variables must be set")
		return // Exit if API keys are not set
	}

	movieService := api.NewMovieService(tmdbAPIKey, omdbAPIKey, appLogger)

	// API endpoint for movie details
	r.HandleFunc("/api/movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}

		// For demonstration, we'll use a placeholder title. In a real app, you might get this from a search or another API call.
		combinedData, err := movieService.GetMovieDetails(id, "") 
		if err != nil {
			appLogger.Error("Error fetching movie details for ID %d: %v", id, err)
			http.Error(w, fmt.Sprintf("Error fetching movie details: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedData)
		appLogger.Success("Successfully fetched movie details for ID %d", id)
	}).Methods("GET")

	// API endpoint for trending movies
	r.HandleFunc("/api/trending", func(w http.ResponseWriter, r *http.Request) {
		trendingMovies, err := movieService.GetTrendingMovies()
		if err != nil {
			appLogger.Error("Error fetching trending movies: %v", err)
			http.Error(w, fmt.Sprintf("Error fetching trending movies: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trendingMovies)
		appLogger.Success("Successfully fetched trending movies")
	}).Methods("GET")

	// API endpoint for movie credits (cast and crew)
	r.HandleFunc("/api/movie/{id}/credits", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}

		credits, err := movieService.GetMovieCredits(id)
		if err != nil {
			appLogger.Error("Error fetching movie credits for ID %d: %v", id, err)
			http.Error(w, fmt.Sprintf("Error fetching movie credits: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(credits)
		appLogger.Success("Successfully fetched movie credits for ID %d", id)
	}).Methods("GET")

	// API endpoint for genres
	r.HandleFunc("/api/genres", func(w http.ResponseWriter, r *http.Request) {
		genres, err := movieService.GetGenres()
		if err != nil {
			appLogger.Error("Error fetching genres: %v", err)
			http.Error(w, fmt.Sprintf("Error fetching genres: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(genres)
		appLogger.Success("Successfully fetched genres")
	}).Methods("GET")

	// API endpoint for movie search
	r.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Search query is required", http.StatusBadRequest)
			return
		}

		searchResults, err := movieService.SearchMovies(query)
		if err != nil {
			appLogger.Error("Error searching movies with query '%s': %v", query, err)
			http.Error(w, fmt.Sprintf("Error searching movies: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResults)
		appLogger.Success("Successfully searched movies with query '%s'", query)
	}).Methods("GET")

	// API endpoint for movie discovery with filters
	r.HandleFunc("/api/discover", func(w http.ResponseWriter, r *http.Request) {
		genreID := r.URL.Query().Get("genre")
		year := r.URL.Query().Get("year")
		sortBy := r.URL.Query().Get("sort_by")

		movies, err := movieService.DiscoverMovies(genreID, year, sortBy)
		if err != nil {
			appLogger.Error("Error discovering movies: %v", err)
			http.Error(w, fmt.Sprintf("Error discovering movies: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies)
		appLogger.Success("Successfully discovered movies with filters")
	}).Methods("GET")

	fmt.Println("Server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
