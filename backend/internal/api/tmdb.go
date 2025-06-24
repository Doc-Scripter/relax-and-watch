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

// GetTrendingContent fetches trending movies or TV shows from TMDB with pagination.
func (c *TMDBClient) GetTrendingContent(contentType string, page int) (map[string]interface{}, error) {
	if contentType == "" {
		contentType = "movie"
	}
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/trending/%s/week?api_key=%s&page=%d", TMDB_BASE_URL, contentType, c.APIKey, page)
	return c.fetchData(url)
}

// GetTVDetails fetches TV show details from TMDB.
func (c *TMDBClient) GetTVDetails(tvID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/tv/%d?api_key=%s", TMDB_BASE_URL, tvID, c.APIKey)
	return c.fetchData(url)
}

// GetTVGenres fetches the list of TV genres from TMDB.
func (c *TMDBClient) GetTVGenres() ([]interface{}, error) {
	url := fmt.Sprintf("%s/genre/tv/list?api_key=%s", TMDB_BASE_URL, c.APIKey)
	data, err := c.fetchData(url)
	if err != nil {
		return nil, err
	}

	if genres, ok := data["genres"].([]interface{}); ok {
		return genres, nil
	}
	return nil, fmt.Errorf("could not find 'genres' array in TMDB TV genres response")
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

// SearchContent searches for movies or TV shows by title from TMDB with pagination.
func (c *TMDBClient) SearchContent(query, contentType string, page int) (map[string]interface{}, error) {
	if contentType == "" {
		contentType = "movie"
	}
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/search/%s?api_key=%s&query=%s&page=%d", TMDB_BASE_URL, contentType, c.APIKey, query, page)
	return c.fetchData(url)
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

// DiscoverContent discovers movies or TV shows with filters from TMDB with pagination.
func (c *TMDBClient) DiscoverContent(contentType string, filters map[string]string, page int) (map[string]interface{}, error) {
	if contentType == "" {
		contentType = "movie"
	}
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/discover/%s?api_key=%s&page=%d", TMDB_BASE_URL, contentType, c.APIKey, page)
	
	// Add filters
	if genreID, ok := filters["genre"]; ok && genreID != "" && genreID != "all" {
		url += fmt.Sprintf("&with_genres=%s", genreID)
	}
	
	if year, ok := filters["year"]; ok && year != "" && year != "all" {
		if contentType == "movie" {
			url += fmt.Sprintf("&year=%s", year)
		} else {
			url += fmt.Sprintf("&first_air_date_year=%s", year)
		}
	}
	
	if rating, ok := filters["rating"]; ok && rating != "" && rating != "all" {
		url += fmt.Sprintf("&vote_average.gte=%s", rating)
	}
	
	if runtime, ok := filters["runtime"]; ok && runtime != "" && runtime != "all" {
		// Parse runtime filter (e.g., "90-120", "180-", "0-90")
		switch runtime {
		case "0-90":
			url += "&with_runtime.lte=90"
		case "90-120":
			url += "&with_runtime.gte=90&with_runtime.lte=120"
		case "120-180":
			url += "&with_runtime.gte=120&with_runtime.lte=180"
		case "180-":
			url += "&with_runtime.gte=180"
		}
	}
	
	if sortBy, ok := filters["sort_by"]; ok && sortBy != "" {
		url += fmt.Sprintf("&sort_by=%s", sortBy)
	} else {
		url += "&sort_by=popularity.desc"
	}

	return c.fetchData(url)
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