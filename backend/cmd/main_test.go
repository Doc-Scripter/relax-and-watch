package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"r.a.w/backend/internal/api"
	"r.a.w/backend/pkg/logger"
)

func TestMain(m *testing.M) {
	// Load environment variables from .env file for tests
	currentDir, _ := os.Getwd()
	envFilePath := filepath.Join(currentDir, "..", "..", ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Printf("Warning: Could not load .env file for tests: %v", err)
	}

	os.Exit(m.Run())
}

func setupRouter() *mux.Router {
	// Initialize custom logger for tests
	currentDir, _ := os.Getwd()
	logFilePath := filepath.Join(currentDir, "..", "logs", "backend_test_errors.log")
	appLogger, err := logger.NewLogger(logFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize logger for tests: %v", err)
	}

	tmdbAPIKey := os.Getenv("TMDB_API_KEY")
	omdbAPIKey := os.Getenv("OMDB_API_KEY")

	if tmdbAPIKey == "" || omdbAPIKey == "" {
		log.Fatalf("TMDB_API_KEY and OMDB_API_KEY environment variables must be set for tests")
	}

	movieService := api.NewMovieService(tmdbAPIKey, omdbAPIKey, appLogger)

	r := mux.NewRouter()

	r.HandleFunc("/api/movie/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}

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

	return r
}

func TestGetMovieDetails(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/movie/550", nil) // Using a known movie ID (Fight Club)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var combinedData struct {
		TMDBData map[string]interface{}
		OMDBData map[string]interface{}
	}
	err := json.Unmarshal(rr.Body.Bytes(), &combinedData)
	assert.NoError(t, err)

	// Assert that at least one of the data sources has a title/plot
	assert.True(t, combinedData.TMDBData != nil || combinedData.OMDBData != nil, "Combined data should not be empty")

	if combinedData.TMDBData != nil {
		assert.NotNil(t, combinedData.TMDBData["title"], "TMDB title should not be nil")
		assert.NotNil(t, combinedData.TMDBData["overview"], "TMDB overview (plot) should not be nil")
	} else if combinedData.OMDBData != nil {
		assert.NotNil(t, combinedData.OMDBData["Title"], "OMDB title should not be nil")
		assert.NotNil(t, combinedData.OMDBData["Plot"], "OMDB plot should not be nil")
	}
}

func TestGetMovieDetailsInvalidID(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/movie/abc", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid movie ID")
}

func TestGetTrendingMovies(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/trending", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var trendingMovies []interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &trendingMovies)
	assert.NoError(t, err)
	assert.NotEmpty(t, trendingMovies, "Trending movies list should not be empty")
}