package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	TMDB_BASE_URL = "https://api.themoviedb.org/3"
)

// TMDBClient represents a client for the TMDB API.
type TMDBClient struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewTMDBClient creates a new TMDB API client.
func NewTMDBClient(apiKey string) *TMDBClient {
	return &TMDBClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetMovieDetails fetches movie details from TMDB.
func (c *TMDBClient) GetMovieDetails(movieID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s", TMDB_BASE_URL, movieID, c.APIKey)
	return c.fetchData(url)
}

// GetTrendingMovies fetches trending movies from TMDB.
func (c *TMDBClient) GetTrendingMovies() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/trending/movie/week?api_key=%s", TMDB_BASE_URL, c.APIKey)
	return c.fetchData(url)
}

// fetchData makes an HTTP GET request and unmarshals the JSON response.
func (c *TMDBClient) fetchData(url string) (map[string]interface{}, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data);
	err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return data, nil
}