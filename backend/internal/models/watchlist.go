package models

import (
	"time"
)

// WatchlistItem represents a single item in a user's watchlist
type WatchlistItem struct {
	ID          string    `json:"id"`
	MovieID     int       `json:"movie_id"`
	Title       string    `json:"title"`
	PosterPath  string    `json:"poster_path"`
	ReleaseDate string    `json:"release_date"`
	Genre       string    `json:"genre"`
	Rating      float64   `json:"rating"`
	Overview    string    `json:"overview"`
	IsWatched   bool      `json:"is_watched"`
	AddedAt     time.Time `json:"added_at"`
	WatchedAt   *time.Time `json:"watched_at,omitempty"`
	UserNotes   string    `json:"user_notes"`
}

// Watchlist represents a user's complete watchlist
type Watchlist struct {
	UserID    string          `json:"user_id"`
	Items     []WatchlistItem `json:"items"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// WatchlistStats represents statistics about a user's watchlist
type WatchlistStats struct {
	TotalItems    int     `json:"total_items"`
	WatchedItems  int     `json:"watched_items"`
	UnwatchedItems int    `json:"unwatched_items"`
	AverageRating float64 `json:"average_rating"`
	TopGenres     []GenreCount `json:"top_genres"`
}

// GenreCount represents genre statistics
type GenreCount struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

// ShareableWatchlist represents a watchlist that can be shared
type ShareableWatchlist struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Items       []WatchlistItem `json:"items"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	IsPublic    bool            `json:"is_public"`
	ShareToken  string          `json:"share_token"`
}