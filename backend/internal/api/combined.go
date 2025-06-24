package api

import (
	"fmt"
	"strings"

	"r.a.w/backend/pkg/logger"
)

// CombinedMovieData represents the combined data from TMDB and OMDB.
type CombinedMovieData struct {
	TMDBData map[string]interface{}
	OMDBData map[string]interface{}
}

// MovieService provides methods to interact with movie APIs.
type MovieService struct {
	TMDBClient *TMDBClient
	OMDBClient *OMDBClient
	Logger     *logger.Logger
}

// NewMovieService creates a new MovieService.
func NewMovieService(tmdbAPIKey, omdbAPIKey string, appLogger *logger.Logger) *MovieService {
	return &MovieService{
		TMDBClient: NewTMDBClient(tmdbAPIKey, appLogger),
		OMDBClient: NewOMDBClient(omdbAPIKey, appLogger),
		Logger:     appLogger,
	}
}

// GetMovieDetails fetches movie details, combining data from TMDB and OMDB.
// It prioritizes TMDB and uses OMDB as a fallback for additional data.
func (s *MovieService) GetMovieDetails(tmdbMovieID int, movieTitle string) (*CombinedMovieData, error) {
	combinedData := &CombinedMovieData{}

	// 1. Fetch from TMDB
	tmdbData, err := s.TMDBClient.GetMovieDetails(tmdbMovieID)
	if err != nil {
		s.Logger.Warning("Error fetching from TMDB for ID %d: %v", tmdbMovieID, err)
		// Continue to OMDB even if TMDB fails, as OMDB might have some data
	} else {
		combinedData.TMDBData = tmdbData
	}

	// 2. Fetch from OMDB using title or IMDB ID from TMDB data
	omdbSearchTitle := movieTitle
	if tmdbData != nil {
		if imdbID, ok := tmdbData["imdb_id"].(string); ok && imdbID != "" {
			omdbData, err := s.OMDBClient.GetMovieByID(imdbID)
			if err != nil {
				s.Logger.Warning("Error fetching from OMDB by IMDB ID %s: %v", imdbID, err)
			} else {
				combinedData.OMDBData = omdbData
			}
		} else if title, ok := tmdbData["title"].(string); ok && title != "" {
			omdbSearchTitle = title // Use TMDB title if available
		}
	}

	// Fallback to searching OMDB by title if no IMDB ID was found or OMDB by ID failed
	if combinedData.OMDBData == nil && omdbSearchTitle != "" {
		omdbData, err := s.OMDBClient.GetMovieByTitle(omdbSearchTitle)
		if err != nil {
			s.Logger.Warning("Error fetching from OMDB by title '%s': %v", omdbSearchTitle, err)
		} else {
			combinedData.OMDBData = omdbData
		}
	}

	// 3. Data Validation (basic example)
	if combinedData.TMDBData == nil && combinedData.OMDBData == nil {
		return nil, fmt.Errorf("could not retrieve movie details from either TMDB or OMDB")
	}

	if err := s.validateMovieData(combinedData); err != nil {
		s.Logger.Warning("Validation warning for movie ID %d: %v", tmdbMovieID, err)
		// Optionally, you can return the data with a warning or return an error
	}

	return combinedData, nil
}

// GetTrendingMovies fetches trending movies from TMDB.
func (s *MovieService) GetTrendingMovies() ([]interface{}, error) {
	return s.TMDBClient.GetTrendingMovies()
}

// validateMovieData performs basic validation on the combined movie data.
func (s *MovieService) validateMovieData(data *CombinedMovieData) error {
	var errors []string

	if data.TMDBData != nil {
		if _, ok := data.TMDBData["title"]; !ok {
			errors = append(errors, "TMDB data missing 'title'")
		}
		if _, ok := data.TMDBData["overview"]; !ok {
			errors = append(errors, "TMDB data missing 'overview'")
		}
	}

	if data.OMDBData != nil {
		if _, ok := data.OMDBData["Title"]; !ok {
			errors = append(errors, "OMDB data missing 'Title'")
		}
		if _, ok := data.OMDBData["imdbRating"]; !ok {
			errors = append(errors, "OMDB data missing 'imdbRating'")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("data validation warnings: %s", strings.Join(errors, "; "))
	}

	return nil
}
