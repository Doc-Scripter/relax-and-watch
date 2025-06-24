package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"r.a.w/backend/internal/handlers"
)

// SetupRoutes configures all the application routes
func SetupRoutes(movieHandler *handlers.MovieHandler, watchlistHandler *handlers.WatchlistHandler) *mux.Router {
	r := mux.NewRouter()

	// Serve static files from the frontend/public directory with proper MIME types
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ext := strings.ToLower(filepath.Ext(path))
		
		// Set proper MIME types
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		
		// Serve the file
		http.FileServer(http.Dir("../../frontend/public")).ServeHTTP(w, r)
	})))
	
	// Serve index.html for the root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/public/index.html")
	})
	
	// Handle favicon.ico requests
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		// Return empty response to prevent 404
		w.WriteHeader(http.StatusNoContent)
	})

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Movie routes
	api.HandleFunc("/movie/{id}", movieHandler.GetMovieDetails).Methods("GET")
	api.HandleFunc("/movie/{id}/credits", movieHandler.GetMovieCredits).Methods("GET")
	api.HandleFunc("/trending", movieHandler.GetTrendingMovies).Methods("GET")
	api.HandleFunc("/genres", movieHandler.GetGenres).Methods("GET")
	api.HandleFunc("/search", movieHandler.SearchMovies).Methods("GET")
	api.HandleFunc("/discover", movieHandler.DiscoverMovies).Methods("GET")
	
	// Watchlist routes
	api.HandleFunc("/watchlist/{userID}", watchlistHandler.GetWatchlist).Methods("GET")
	api.HandleFunc("/watchlist/{userID}", watchlistHandler.AddToWatchlist).Methods("POST")
	api.HandleFunc("/watchlist/{userID}/{itemID}", watchlistHandler.RemoveFromWatchlist).Methods("DELETE")
	api.HandleFunc("/watchlist/{userID}/{itemID}/watched", watchlistHandler.MarkAsWatched).Methods("PUT")
	api.HandleFunc("/watchlist/{userID}/{itemID}/unwatched", watchlistHandler.MarkAsUnwatched).Methods("PUT")
	api.HandleFunc("/watchlist/{userID}/stats", watchlistHandler.GetWatchlistStats).Methods("GET")
	api.HandleFunc("/watchlist/{userID}/export", watchlistHandler.ExportWatchlist).Methods("GET")
	api.HandleFunc("/watchlist/{userID}/share", watchlistHandler.CreateShareableWatchlist).Methods("POST")
	api.HandleFunc("/shared/{shareToken}", watchlistHandler.GetSharedWatchlist).Methods("GET")

	return r
}