package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"r.a.w/backend/pkg/logger"
)

const (
	TMDB_BASE_URL = "https://api.themoviedb.org/3"
)

// TMDBClient represents a client for the TMDB API.
type TMDBClient struct {
	APIKey     string
	HTTPClient *http.Client
	Logger     *logger.Logger
}

// NewTMDBClient creates a new TMDB API client.
func NewTMDBClient(apiKey string, appLogger *logger.Logger) *TMDBClient {
	return &TMDBClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
		Logger: appLogger,
	}
}

// GetMovieDetails fetches movie details from TMDB.
func (c *TMDBClient) GetMovieDetails(movieID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s", TMDB_BASE_URL, movieID, c.APIKey)
	return c.fetchData(url)
}

// GetTrendingMovies fetches trending movies from TMDB.
func (c *TMDBClient) GetTrendingMovies() ([]interface{}, error) {
	url := fmt.Sprintf("%s/trending/movie/week?api_key=%s", TMDB_BASE_URL, c.APIKey)
	data, err := c.fetchData(url)
	if err != nil {
		return nil, err
	}

	// TMDB trending API returns a map with a "results" key containing the array of movies
	if results, ok := data["results"].([]interface{}); ok {
		return results, nil
	}
	return nil, fmt.Errorf("could not find 'results' array in TMDB trending response")
}

// GetMovieCredits fetches cast and crew information for a movie from TMDB.
func (c *TMDBClient) GetMovieCredits(movieID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%d/credits?api_key=%s", TMDB_BASE_URL, movieID, c.APIKey)
	return c.fetchData(url)
}

// GetGenres fetches the list of movie genres from TMDB.
func (c *TMDBClient) GetGenres() ([]interface{}, error) {
	url := fmt.Sprintf("%s/genre/movie/list?api_key=%s", TMDB_BASE_URL, c.APIKey)
	data, err := c.fetchData(url)
	if err != nil {
		return nil, err
	}

	if genres, ok := data["genres"].([]interface{}); ok {
		return genres, nil
	}
	return nil, fmt.Errorf("could not find 'genres' array in TMDB genres response")
}

// SearchMovies searches for movies by title from TMDB.
func (c *TMDBClient) SearchMovies(query string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", TMDB_BASE_URL, c.APIKey, query)
	data, err := c.fetchData(url)
	if err != nil {
		return nil, err
	}

	if results, ok := data["results"].([]interface{}); ok {
		return results, nil
	}
	return nil, fmt.Errorf("could not find 'results' array in TMDB search response")
}

// DiscoverMovies discovers movies with filters from TMDB.
func (c *TMDBClient) DiscoverMovies(genreID, year, sortBy string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/discover/movie?api_key=%s", TMDB_BASE_URL, c.APIKey)
	
	if genreID != "" && genreID != "all" {
		url += fmt.Sprintf("&with_genres=%s", genreID)
	}
	if year != "" && year != "all" {
		url += fmt.Sprintf("&year=%s", year)
	}
	if sortBy != "" {
		url += fmt.Sprintf("&sort_by=%s", sortBy)
	} else {
		url += "&sort_by=popularity.desc"
	}

	data, err := c.fetchData(url)
	if err != nil {
		return nil, err
	}

	if results, ok := data["results"].([]interface{}); ok {
		return results, nil
	}
	return nil, fmt.Errorf("could not find 'results' array in TMDB discover response")
}

// fetchData makes an HTTP GET request and unmarshals the JSON response.
func (c *TMDBClient) fetchData(url string) (map[string]interface{}, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		c.Logger.Error("Failed to make request to TMDB: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Logger.Error("TMDB API request failed with status code: %d for URL: %s", resp.StatusCode, url)
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Error("Failed to read response body from TMDB: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.Logger.Error("Failed to unmarshal JSON response from TMDB: %v", err)
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return data, nil
}