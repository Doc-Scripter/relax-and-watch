package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"r.a.w/backend/internal/api"
	"r.a.w/backend/pkg/logger"
)

// MovieHandler handles movie-related HTTP requests
type MovieHandler struct {
	MovieService *api.MovieService
	Logger       *logger.Logger
}

// NewMovieHandler creates a new MovieHandler
func NewMovieHandler(movieService *api.MovieService, logger *logger.Logger) *MovieHandler {
	return &MovieHandler{
		MovieService: movieService,
		Logger:       logger,
	}
}

// GetMovieDetails handles GET /api/movie/{id}
func (h *MovieHandler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// First get basic movie details from TMDB to extract the title
	tmdbData, err := h.MovieService.TMDBClient.GetMovieDetails(id)
	if err != nil {
		h.Logger.Error("Error fetching TMDB movie details for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Error fetching movie details: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract title from TMDB data
	movieTitle := ""
	if title, ok := tmdbData["title"].(string); ok {
		movieTitle = title
	}

	// Now get combined data with the proper title
	combinedData, err := h.MovieService.GetMovieDetails(id, movieTitle)
	if err != nil {
		h.Logger.Error("Error fetching combined movie details for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Error fetching movie details: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(combinedData)
	h.Logger.Success("Successfully fetched movie details for ID %d", id)
}

// GetTrendingMovies handles GET /api/trending
func (h *MovieHandler) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	trendingMovies, err := h.MovieService.GetTrendingMovies()
	if err != nil {
		h.Logger.Error("Error fetching trending movies: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching trending movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trendingMovies)
	h.Logger.Success("Successfully fetched trending movies")
}

// GetMovieCredits handles GET /api/movie/{id}/credits
func (h *MovieHandler) GetMovieCredits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	credits, err := h.MovieService.GetMovieCredits(id)
	if err != nil {
		h.Logger.Error("Error fetching movie credits for ID %d: %v", id, err)
		http.Error(w, fmt.Sprintf("Error fetching movie credits: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credits)
	h.Logger.Success("Successfully fetched movie credits for ID %d", id)
}

// GetGenres handles GET /api/genres
func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.MovieService.GetGenres()
	if err != nil {
		h.Logger.Error("Error fetching genres: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching genres: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
	h.Logger.Success("Successfully fetched genres")
}

// SearchMovies handles GET /api/search
func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	searchResults, err := h.MovieService.SearchMovies(query)
	if err != nil {
		h.Logger.Error("Error searching movies with query '%s': %v", query, err)
		http.Error(w, fmt.Sprintf("Error searching movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
	h.Logger.Success("Successfully searched movies with query '%s'", query)
}

// DiscoverMovies handles GET /api/discover
func (h *MovieHandler) DiscoverMovies(w http.ResponseWriter, r *http.Request) {
	genreID := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")
	sortBy := r.URL.Query().Get("sort_by")

	movies, err := h.MovieService.DiscoverMovies(genreID, year, sortBy)
	if err != nil {
		h.Logger.Error("Error discovering movies: %v", err)
		http.Error(w, fmt.Sprintf("Error discovering movies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	h.Logger.Success("Successfully discovered movies with filters")
}