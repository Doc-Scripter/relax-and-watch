package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"r.a.w/backend/internal/models"
	"r.a.w/backend/pkg/logger"
)

// ExportService handles exporting watchlists to various formats
type ExportService struct {
	logger *logger.Logger
}

// NewExportService creates a new export service
func NewExportService(logger *logger.Logger) *ExportService {
	return &ExportService{
		logger: logger,
	}
}

// ExportToCSV exports a watchlist to CSV format
func (s *ExportService) ExportToCSV(watchlist *models.Watchlist) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// Write header
	header := []string{
		"Title",
		"Release Date",
		"Genre",
		"Rating",
		"Status",
		"Added Date",
		"Watched Date",
		"Notes",
		"Overview",
	}
	
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}
	
	// Write data rows
	for _, item := range watchlist.Items {
		status := "Unwatched"
		watchedDate := ""
		
		if item.IsWatched {
			status = "Watched"
			if item.WatchedAt != nil {
				watchedDate = item.WatchedAt.Format("2006-01-02")
			}
		}
		
		row := []string{
			item.Title,
			item.ReleaseDate,
			item.Genre,
			fmt.Sprintf("%.1f", item.Rating),
			status,
			item.AddedAt.Format("2006-01-02"),
			watchedDate,
			item.UserNotes,
			item.Overview,
		}
		
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}
	
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}
	
	s.logger.Success("Watchlist exported to CSV for user %s", watchlist.UserID)
	return buf.Bytes(), nil
}

// ExportToPDF exports a watchlist to PDF format (simplified HTML-based approach)
func (s *ExportService) ExportToPDF(watchlist *models.Watchlist, stats *models.WatchlistStats) ([]byte, error) {
	// Generate HTML content that can be converted to PDF
	html := s.generateHTMLReport(watchlist, stats)
	
	// For now, return HTML bytes. In a production environment, you would use a library like wkhtmltopdf
	// or a Go PDF library to convert HTML to PDF
	s.logger.Success("Watchlist exported to HTML/PDF for user %s", watchlist.UserID)
	return []byte(html), nil
}

// generateHTMLReport generates an HTML report of the watchlist
func (s *ExportService) generateHTMLReport(watchlist *models.Watchlist, stats *models.WatchlistStats) string {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>My Watchlist Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; color: #333; }
        .header { text-align: center; margin-bottom: 30px; }
        .stats { background: #f5f5f5; padding: 20px; border-radius: 8px; margin-bottom: 30px; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; }
        .stat-item { text-align: center; }
        .stat-number { font-size: 24px; font-weight: bold; color: #e94560; }
        .stat-label { font-size: 14px; color: #666; }
        .movies-section { margin-top: 30px; }
        .movie-item { border-bottom: 1px solid #eee; padding: 15px 0; display: flex; align-items: flex-start; }
        .movie-info { flex: 1; }
        .movie-title { font-size: 18px; font-weight: bold; margin-bottom: 5px; }
        .movie-details { color: #666; font-size: 14px; margin-bottom: 5px; }
        .movie-overview { color: #888; font-size: 13px; line-height: 1.4; }
        .status-watched { color: #28a745; font-weight: bold; }
        .status-unwatched { color: #ffc107; font-weight: bold; }
        .rating { color: #e94560; font-weight: bold; }
        .genres { margin-top: 20px; }
        .genre-item { display: inline-block; background: #e94560; color: white; padding: 5px 10px; margin: 2px; border-radius: 15px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>My Watchlist Report</h1>
        <p>Generated on ` + time.Now().Format("January 2, 2006") + `</p>
    </div>
    
    <div class="stats">
        <h2>Statistics</h2>
        <div class="stats-grid">
            <div class="stat-item">
                <div class="stat-number">` + strconv.Itoa(stats.TotalItems) + `</div>
                <div class="stat-label">Total Movies</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">` + strconv.Itoa(stats.WatchedItems) + `</div>
                <div class="stat-label">Watched</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">` + strconv.Itoa(stats.UnwatchedItems) + `</div>
                <div class="stat-label">To Watch</div>
            </div>
            <div class="stat-item">
                <div class="stat-number">` + fmt.Sprintf("%.1f", stats.AverageRating) + `</div>
                <div class="stat-label">Avg Rating</div>
            </div>
        </div>`
	
	if len(stats.TopGenres) > 0 {
		html += `
        <div class="genres">
            <h3>Top Genres</h3>`
		for _, genre := range stats.TopGenres {
			html += `<span class="genre-item">` + genre.Genre + ` (` + strconv.Itoa(genre.Count) + `)</span>`
		}
		html += `
        </div>`
	}
	
	html += `
    </div>
    
    <div class="movies-section">
        <h2>Movies (` + strconv.Itoa(len(watchlist.Items)) + `)</h2>`
	
	for _, item := range watchlist.Items {
		status := `<span class="status-unwatched">To Watch</span>`
		watchedInfo := ""
		
		if item.IsWatched {
			status = `<span class="status-watched">Watched</span>`
			if item.WatchedAt != nil {
				watchedInfo = " on " + item.WatchedAt.Format("Jan 2, 2006")
			}
		}
		
		html += `
        <div class="movie-item">
            <div class="movie-info">
                <div class="movie-title">` + item.Title + `</div>
                <div class="movie-details">
                    ` + item.ReleaseDate + ` • ` + item.Genre + ` • <span class="rating">★ ` + fmt.Sprintf("%.1f", item.Rating) + `</span> • ` + status + watchedInfo + `
                </div>`
		
		if item.UserNotes != "" {
			html += `<div class="movie-details"><strong>Notes:</strong> ` + item.UserNotes + `</div>`
		}
		
		html += `
                <div class="movie-overview">` + item.Overview + `</div>
            </div>
        </div>`
	}
	
	html += `
    </div>
</body>
</html>`
	
	return html
}