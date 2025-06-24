package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"r.a.w/backend/internal/models"
	"r.a.w/backend/internal/services"
	"r.a.w/backend/pkg/logger"
)

// WatchlistHandler handles watchlist-related HTTP requests
type WatchlistHandler struct {
	WatchlistService *services.WatchlistService
	ExportService    *services.ExportService
	Logger           *logger.Logger
}

// NewWatchlistHandler creates a new WatchlistHandler
func NewWatchlistHandler(watchlistService *services.WatchlistService, exportService *services.ExportService, logger *logger.Logger) *WatchlistHandler {
	return &WatchlistHandler{
		WatchlistService: watchlistService,
		ExportService:    exportService,
		Logger:           logger,
	}
}

// GetWatchlist handles GET /api/watchlist/{userID}
func (h *WatchlistHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	watchlist, err := h.WatchlistService.GetWatchlist(userID)
	if err != nil {
		h.Logger.Error("Error fetching watchlist for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error fetching watchlist: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(watchlist)
	h.Logger.Success("Successfully fetched watchlist for user %s", userID)
}

// AddToWatchlist handles POST /api/watchlist/{userID}
func (h *WatchlistHandler) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	var item models.WatchlistItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.WatchlistService.AddToWatchlist(userID, item); err != nil {
		h.Logger.Error("Error adding to watchlist for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error adding to watchlist: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Movie added to watchlist"})
	h.Logger.Success("Successfully added movie to watchlist for user %s", userID)
}

// RemoveFromWatchlist handles DELETE /api/watchlist/{userID}/{itemID}
func (h *WatchlistHandler) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	itemID := vars["itemID"]
	
	if userID == "" || itemID == "" {
		http.Error(w, "User ID and Item ID are required", http.StatusBadRequest)
		return
	}
	
	if err := h.WatchlistService.RemoveFromWatchlist(userID, itemID); err != nil {
		h.Logger.Error("Error removing from watchlist for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error removing from watchlist: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Movie removed from watchlist"})
	h.Logger.Success("Successfully removed movie from watchlist for user %s", userID)
}

// MarkAsWatched handles PUT /api/watchlist/{userID}/{itemID}/watched
func (h *WatchlistHandler) MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	itemID := vars["itemID"]
	
	if userID == "" || itemID == "" {
		http.Error(w, "User ID and Item ID are required", http.StatusBadRequest)
		return
	}
	
	var requestBody struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)
	
	if err := h.WatchlistService.MarkAsWatched(userID, itemID, requestBody.Notes); err != nil {
		h.Logger.Error("Error marking as watched for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error marking as watched: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Movie marked as watched"})
	h.Logger.Success("Successfully marked movie as watched for user %s", userID)
}

// MarkAsUnwatched handles PUT /api/watchlist/{userID}/{itemID}/unwatched
func (h *WatchlistHandler) MarkAsUnwatched(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	itemID := vars["itemID"]
	
	if userID == "" || itemID == "" {
		http.Error(w, "User ID and Item ID are required", http.StatusBadRequest)
		return
	}
	
	if err := h.WatchlistService.MarkAsUnwatched(userID, itemID); err != nil {
		h.Logger.Error("Error marking as unwatched for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error marking as unwatched: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Movie marked as unwatched"})
	h.Logger.Success("Successfully marked movie as unwatched for user %s", userID)
}

// GetWatchlistStats handles GET /api/watchlist/{userID}/stats
func (h *WatchlistHandler) GetWatchlistStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	stats, err := h.WatchlistService.GetWatchlistStats(userID)
	if err != nil {
		h.Logger.Error("Error fetching watchlist stats for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error fetching watchlist stats: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
	h.Logger.Success("Successfully fetched watchlist stats for user %s", userID)
}

// ExportWatchlist handles GET /api/watchlist/{userID}/export
func (h *WatchlistHandler) ExportWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	format := r.URL.Query().Get("format") // csv or pdf
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	if format == "" {
		format = "csv"
	}
	
	watchlist, err := h.WatchlistService.GetWatchlist(userID)
	if err != nil {
		h.Logger.Error("Error fetching watchlist for export for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error fetching watchlist: %v", err), http.StatusInternalServerError)
		return
	}
	
	switch format {
	case "csv":
		data, err := h.ExportService.ExportToCSV(watchlist)
		if err != nil {
			h.Logger.Error("Error exporting to CSV for user %s: %v", userID, err)
			http.Error(w, fmt.Sprintf("Error exporting to CSV: %v", err), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=watchlist_%s.csv", userID))
		w.Write(data)
		
	case "pdf":
		stats, err := h.WatchlistService.GetWatchlistStats(userID)
		if err != nil {
			h.Logger.Error("Error fetching stats for PDF export for user %s: %v", userID, err)
			http.Error(w, fmt.Sprintf("Error fetching stats: %v", err), http.StatusInternalServerError)
			return
		}
		
		data, err := h.ExportService.ExportToPDF(watchlist, stats)
		if err != nil {
			h.Logger.Error("Error exporting to PDF for user %s: %v", userID, err)
			http.Error(w, fmt.Sprintf("Error exporting to PDF: %v", err), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=watchlist_%s.html", userID))
		w.Write(data)
		
	default:
		http.Error(w, "Invalid format. Use 'csv' or 'pdf'", http.StatusBadRequest)
		return
	}
	
	h.Logger.Success("Successfully exported watchlist in %s format for user %s", format, userID)
}

// CreateShareableWatchlist handles POST /api/watchlist/{userID}/share
func (h *WatchlistHandler) CreateShareableWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	
	var requestBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	shareableWatchlist, err := h.WatchlistService.CreateShareableWatchlist(userID, requestBody.Title, requestBody.Description, requestBody.IsPublic)
	if err != nil {
		h.Logger.Error("Error creating shareable watchlist for user %s: %v", userID, err)
		http.Error(w, fmt.Sprintf("Error creating shareable watchlist: %v", err), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shareableWatchlist)
	h.Logger.Success("Successfully created shareable watchlist for user %s", userID)
}

// GetSharedWatchlist handles GET /api/shared/{shareToken}
func (h *WatchlistHandler) GetSharedWatchlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shareToken := vars["shareToken"]
	
	if shareToken == "" {
		http.Error(w, "Share token is required", http.StatusBadRequest)
		return
	}
	
	sharedWatchlist, err := h.WatchlistService.GetSharedWatchlist(shareToken)
	if err != nil {
		h.Logger.Error("Error fetching shared watchlist with token %s: %v", shareToken, err)
		http.Error(w, fmt.Sprintf("Error fetching shared watchlist: %v", err), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sharedWatchlist)
	h.Logger.Success("Successfully fetched shared watchlist with token %s", shareToken)
}