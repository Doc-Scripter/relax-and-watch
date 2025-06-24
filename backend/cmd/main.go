package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"r.a.w/backend/internal/api"
)

// Serve static files from the "frontend/public" directory
func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("../../frontend/public"))
	r.PathPrefix("/").Handler(fs)

	// Initialize API clients
	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	omdbAPIKey := os.Getenv("OMDB_API_KEY")

	if tmdbAPIKey == "" || omdbAPIKey == "" {
		log.Fatal("TMDB_API_KEY and OMDB_API_KEY environment variables must be set")
	}

	movieService := api.NewMovieService(tmdbAPIKey, omdbAPIKey)

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
			http.Error(w, fmt.Sprintf("Error fetching movie details: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedData)
	}).Methods("GET")

	// API endpoint for trending movies
	r.HandleFunc("/api/trending", func(w http.ResponseWriter, r *http.Request) {
		trendingMovies, err := movieService.GetTrendingMovies()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching trending movies: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trendingMovies)
	}).Methods("GET")

	fmt.Println("Server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
