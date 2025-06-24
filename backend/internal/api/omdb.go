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
	OMDB_BASE_URL = "http://www.omdbapi.com/"
)

// OMDBClient represents a client for the OMDB API.
type OMDBClient struct {
	APIKey     string
	HTTPClient *http.Client
	Logger     *logger.Logger
}

// NewOMDBClient creates a new OMDB API client.
func NewOMDBClient(apiKey string, appLogger *logger.Logger) *OMDBClient {
	return &OMDBClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
		Logger: appLogger,
	}
}

// GetMovieByTitle fetches movie details from OMDB by title.
func (c *OMDBClient) GetMovieByTitle(title string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?t=%s&apikey=%s", OMDB_BASE_URL, title, c.APIKey)
	return c.fetchData(url)
}

// GetMovieByID fetches movie details from OMDB by IMDB ID.
func (c *OMDBClient) GetMovieByID(imdbID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?i=%s&apikey=%s", OMDB_BASE_URL, imdbID, c.APIKey)
	return c.fetchData(url)
}

// fetchData makes an HTTP GET request and unmarshals the JSON response.
func (c *OMDBClient) fetchData(url string) (map[string]interface{}, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		c.Logger.Error("Failed to make request to OMDB: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Logger.Error("OMDB API request failed with status code: %d for URL: %s", resp.StatusCode, url)
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Error("Failed to read response body from OMDB: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.Logger.Error("Failed to unmarshal JSON response from OMDB: %v", err)
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	// OMDB returns a JSON object with "Response":"False" and an "Error" field if the movie is not found.
	if resp, ok := data["Response"]; ok && resp == "False" {
		if errMsg, ok := data["Error"]; ok {
			c.Logger.Error("OMDB API error: %s", errMsg)
			return nil, fmt.Errorf("OMDB API error: %s", errMsg)
		}
		c.Logger.Error("OMDB API error: movie not found or other issue for URL: %s", url)
		return nil, fmt.Errorf("OMDB API error: movie not found or other issue")
	}

	return data, nil
}
