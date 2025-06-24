package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"r.a.w/backend/internal/handlers"
)

// SetupRoutes configures all the application routes
func SetupRoutes(movieHandler *handlers.MovieHandler) *mux.Router {
	r := mux.NewRouter()

	// Serve static files from the frontend/public directory
	fs := http.FileServer(http.Dir("../../frontend/public"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	// Serve index.html for the root path
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../frontend/public/index.html")
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

	return r
}