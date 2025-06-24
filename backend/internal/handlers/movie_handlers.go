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
		http.Error(w, "Invalid content ID", http.StatusBadRequest)
		return
	}

	contentType := r.URL.Query().Get("type")
	if contentType == "" {
		contentType = "movie" // Default to movie for backward compatibility
	}

	if contentType == "tv" {
		// Handle TV show details
		tmdbData, err := h.MovieService.TMDBClient.GetTVDetails(id)
		if err != nil {
			h.Logger.Error("Error fetching TMDB TV details for ID %d: %v", id, err)
			http.Error(w, fmt.Sprintf("Error fetching TV details: %v", err), http.StatusInternalServerError)
			return
		}

		// Extract title from TMDB data
		tvTitle := ""
		if name, ok := tmdbData["name"].(string); ok {
			tvTitle = name
		}

		// Get combined data with the proper title
		combinedData, err := h.MovieService.GetTVDetails(id, tvTitle)
		if err != nil {
			h.Logger.Error("Error fetching combined TV details for ID %d: %v", id, err)
			http.Error(w, fmt.Sprintf("Error fetching TV details: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(combinedData)
		h.Logger.Success("Successfully fetched TV details for ID %d", id)
		return
	}

	// Handle movie details (original logic)
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
	contentType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	
	if contentType != "" || pageStr != "" {
		// Use new paginated endpoint
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		if contentType == "" {
			contentType = "movie"
		}
		
		trendingContent, err := h.MovieService.GetTrendingContent(contentType, page)
		if err != nil {
			h.Logger.Error("Error fetching trending content: %v", err)
			http.Error(w, fmt.Sprintf("Error fetching trending content: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trendingContent)
		h.Logger.Success("Successfully fetched trending %s (page %d)", contentType, page)
		return
	}
	
	// Fallback to old endpoint for backward compatibility
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
	contentType := r.URL.Query().Get("type")
	
	if contentType != "" {
		// Use new type-specific endpoint
		genres, err := h.MovieService.GetGenresByType(contentType)
		if err != nil {
			h.Logger.Error("Error fetching %s genres: %v", contentType, err)
			http.Error(w, fmt.Sprintf("Error fetching genres: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(genres)
		h.Logger.Success("Successfully fetched %s genres", contentType)
		return
	}
	
	// Fallback to old endpoint for backward compatibility
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

	contentType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	
	if contentType != "" || pageStr != "" {
		// Use new paginated search endpoint
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		if contentType == "" {
			contentType = "movie"
		}
		
		searchResults, err := h.MovieService.SearchContent(query, contentType, page)
		if err != nil {
			h.Logger.Error("Error searching %s with query '%s': %v", contentType, query, err)
			http.Error(w, fmt.Sprintf("Error searching content: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResults)
		h.Logger.Success("Successfully searched %s with query '%s' (page %d)", contentType, query, page)
		return
	}
	
	// Fallback to old endpoint for backward compatibility
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
	contentType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	
	if contentType != "" || pageStr != "" {
		// Use new paginated discover endpoint
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		if contentType == "" {
			contentType = "movie"
		}
		
		// Collect filters
		filters := make(map[string]string)
		if genre := r.URL.Query().Get("genre"); genre != "" {
			filters["genre"] = genre
		}
		if year := r.URL.Query().Get("year"); year != "" {
			filters["year"] = year
		}
		if rating := r.URL.Query().Get("rating"); rating != "" {
			filters["rating"] = rating
		}
		if runtime := r.URL.Query().Get("runtime"); runtime != "" {
			filters["runtime"] = runtime
		}
		if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
			filters["sort_by"] = sortBy
		}
		
		content, err := h.MovieService.DiscoverContent(contentType, filters, page)
		if err != nil {
			h.Logger.Error("Error discovering %s: %v", contentType, err)
			http.Error(w, fmt.Sprintf("Error discovering content: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(content)
		h.Logger.Success("Successfully discovered %s with filters (page %d)", contentType, page)
		return
	}
	
	// Fallback to old endpoint for backward compatibility
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