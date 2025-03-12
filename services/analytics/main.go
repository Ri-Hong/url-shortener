// This service provides analytics for short URLs.

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type URLAnalytics struct {
	db *sql.DB
}

type AnalyticsResponse struct {
	ShortCode  string `json:"short_code"`
	ClickCount int    `json:"click_count"`
}

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", os.Getenv("LOCAL_DATABASE_URL")+"?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create analytics service
	analytics := &URLAnalytics{db: db}

	// Initialize router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/analytics/{shortCode}", analytics.GetAnalyticsHandler).Methods("GET")

	// Start server
	log.Println("Starting URL analytics service on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func (s *URLAnalytics) GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the short code from the URL path
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	if shortCode == "" {
		http.Error(w, "No short code provided", http.StatusBadRequest)
		return
	}

	// Look up the click count in the database
	var clickCount int
	err := s.db.QueryRow("SELECT COALESCE(click_count, 0) FROM urls WHERE short_code = $1", shortCode).Scan(&clickCount)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Short URL not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Create response
	response := AnalyticsResponse{
		ShortCode:  shortCode,
		ClickCount: clickCount,
	}

	// Set content type and encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Analytics retrieved for %s: %d clicks", shortCode, clickCount)
}
