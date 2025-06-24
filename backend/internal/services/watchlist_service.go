package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"r.a.w/backend/internal/models"
	"r.a.w/backend/pkg/logger"
)

// WatchlistService handles watchlist operations
type WatchlistService struct {
	dataDir string
	logger  *logger.Logger
}

// NewWatchlistService creates a new watchlist service
func NewWatchlistService(dataDir string, logger *logger.Logger) *WatchlistService {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.Error("Failed to create data directory: %v", err)
	}
	
	return &WatchlistService{
		dataDir: dataDir,
		logger:  logger,
	}
}

// GetWatchlist retrieves a user's watchlist
func (s *WatchlistService) GetWatchlist(userID string) (*models.Watchlist, error) {
	filePath := filepath.Join(s.dataDir, fmt.Sprintf("watchlist_%s.json", userID))
	
	// If file doesn't exist, return empty watchlist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &models.Watchlist{
			UserID:    userID,
			Items:     []models.WatchlistItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read watchlist file: %w", err)
	}
	
	var watchlist models.Watchlist
	if err := json.Unmarshal(data, &watchlist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal watchlist: %w", err)
	}
	
	return &watchlist, nil
}

// AddToWatchlist adds a movie to the user's watchlist
func (s *WatchlistService) AddToWatchlist(userID string, item models.WatchlistItem) error {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return err
	}
	
	// Check if item already exists
	for _, existingItem := range watchlist.Items {
		if existingItem.MovieID == item.MovieID {
			return fmt.Errorf("movie already in watchlist")
		}
	}
	
	// Generate unique ID for the item
	item.ID = s.generateID()
	item.AddedAt = time.Now()
	item.IsWatched = false
	
	watchlist.Items = append(watchlist.Items, item)
	watchlist.UpdatedAt = time.Now()
	
	return s.saveWatchlist(watchlist)
}

// RemoveFromWatchlist removes a movie from the user's watchlist
func (s *WatchlistService) RemoveFromWatchlist(userID, itemID string) error {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return err
	}
	
	// Find and remove the item
	for i, item := range watchlist.Items {
		if item.ID == itemID {
			watchlist.Items = append(watchlist.Items[:i], watchlist.Items[i+1:]...)
			watchlist.UpdatedAt = time.Now()
			return s.saveWatchlist(watchlist)
		}
	}
	
	return fmt.Errorf("item not found in watchlist")
}

// MarkAsWatched marks a movie as watched in the user's watchlist
func (s *WatchlistService) MarkAsWatched(userID, itemID string, notes string) error {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return err
	}
	
	// Find and update the item
	for i, item := range watchlist.Items {
		if item.ID == itemID {
			now := time.Now()
			watchlist.Items[i].IsWatched = true
			watchlist.Items[i].WatchedAt = &now
			watchlist.Items[i].UserNotes = notes
			watchlist.UpdatedAt = time.Now()
			return s.saveWatchlist(watchlist)
		}
	}
	
	return fmt.Errorf("item not found in watchlist")
}

// MarkAsUnwatched marks a movie as unwatched in the user's watchlist
func (s *WatchlistService) MarkAsUnwatched(userID, itemID string) error {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return err
	}
	
	// Find and update the item
	for i, item := range watchlist.Items {
		if item.ID == itemID {
			watchlist.Items[i].IsWatched = false
			watchlist.Items[i].WatchedAt = nil
			watchlist.Items[i].UserNotes = ""
			watchlist.UpdatedAt = time.Now()
			return s.saveWatchlist(watchlist)
		}
	}
	
	return fmt.Errorf("item not found in watchlist")
}

// GetWatchlistStats returns statistics about the user's watchlist
func (s *WatchlistService) GetWatchlistStats(userID string) (*models.WatchlistStats, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}
	
	stats := &models.WatchlistStats{
		TotalItems: len(watchlist.Items),
	}
	
	var totalRating float64
	var ratingCount int
	genreCount := make(map[string]int)
	
	for _, item := range watchlist.Items {
		if item.IsWatched {
			stats.WatchedItems++
		} else {
			stats.UnwatchedItems++
		}
		
		if item.Rating > 0 {
			totalRating += item.Rating
			ratingCount++
		}
		
		// Count genres
		genres := strings.Split(item.Genre, ", ")
		for _, genre := range genres {
			genre = strings.TrimSpace(genre)
			if genre != "" {
				genreCount[genre]++
			}
		}
	}
	
	if ratingCount > 0 {
		stats.AverageRating = totalRating / float64(ratingCount)
	}
	
	// Convert genre map to sorted slice
	for genre, count := range genreCount {
		stats.TopGenres = append(stats.TopGenres, models.GenreCount{
			Genre: genre,
			Count: count,
		})
	}
	
	// Sort genres by count (descending)
	sort.Slice(stats.TopGenres, func(i, j int) bool {
		return stats.TopGenres[i].Count > stats.TopGenres[j].Count
	})
	
	// Keep only top 5 genres
	if len(stats.TopGenres) > 5 {
		stats.TopGenres = stats.TopGenres[:5]
	}
	
	return stats, nil
}

// CreateShareableWatchlist creates a shareable version of the watchlist
func (s *WatchlistService) CreateShareableWatchlist(userID, title, description string, isPublic bool) (*models.ShareableWatchlist, error) {
	watchlist, err := s.GetWatchlist(userID)
	if err != nil {
		return nil, err
	}
	
	shareableWatchlist := &models.ShareableWatchlist{
		ID:          s.generateID(),
		Title:       title,
		Description: description,
		Items:       watchlist.Items,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		IsPublic:    isPublic,
		ShareToken:  s.generateShareToken(),
	}
	
	// Save shareable watchlist
	filePath := filepath.Join(s.dataDir, fmt.Sprintf("shared_watchlist_%s.json", shareableWatchlist.ID))
	data, err := json.MarshalIndent(shareableWatchlist, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shareable watchlist: %w", err)
	}
	
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to save shareable watchlist: %w", err)
	}
	
	return shareableWatchlist, nil
}

// GetSharedWatchlist retrieves a shared watchlist by token
func (s *WatchlistService) GetSharedWatchlist(shareToken string) (*models.ShareableWatchlist, error) {
	// Search for the shared watchlist file with the given token
	files, err := ioutil.ReadDir(s.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}
	
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "shared_watchlist_") {
			filePath := filepath.Join(s.dataDir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}
			
			var sharedWatchlist models.ShareableWatchlist
			if err := json.Unmarshal(data, &sharedWatchlist); err != nil {
				continue
			}
			
			if sharedWatchlist.ShareToken == shareToken {
				return &sharedWatchlist, nil
			}
		}
	}
	
	return nil, fmt.Errorf("shared watchlist not found")
}

// saveWatchlist saves the watchlist to file
func (s *WatchlistService) saveWatchlist(watchlist *models.Watchlist) error {
	filePath := filepath.Join(s.dataDir, fmt.Sprintf("watchlist_%s.json", watchlist.UserID))
	
	data, err := json.MarshalIndent(watchlist, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal watchlist: %w", err)
	}
	
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to save watchlist: %w", err)
	}
	
	s.logger.Success("Watchlist saved for user %s", watchlist.UserID)
	return nil
}

// generateID generates a unique ID
func (s *WatchlistService) generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateShareToken generates a unique share token
func (s *WatchlistService) generateShareToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}