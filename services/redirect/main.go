// This service redirects short URLs to their original long URLs.

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type URLRedirector struct {
	db *sql.DB
}

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", os.Getenv("LOCAL_DATABASE_URL")+"?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create redirector service
	redirector := &URLRedirector{db: db}

	// Initialize router
	r := mux.NewRouter()

	// Define routes - the root path handles redirects
	r.HandleFunc("/{shortCode}", redirector.RedirectHandler).Methods("GET")

	// Start server
	log.Println("Starting URL redirector service on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func (s *URLRedirector) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the short code from the URL path
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]
	
	if shortCode == "" {
		http.Error(w, "No short code provided", http.StatusBadRequest)
		return
	}

	// Look up the original URL in the database
	var longURL string
	err := s.db.QueryRow("SELECT long_url FROM urls WHERE short_code = $1", shortCode).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Short URL not found", http.StatusNotFound)
		} else {
			log.Printf("Database error: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// Increment the click count
	_, err = s.db.Exec("UPDATE urls SET click_count = COALESCE(click_count, 0) + 1 WHERE short_code = $1", shortCode)
	if err != nil {
		log.Printf("Failed to update click count: %v", err)
		// Continue with redirect even if update fails
	}

	// Redirect to the original URL
	http.Redirect(w, r, longURL, http.StatusFound)
	log.Printf("Redirected %s to %s", shortCode, longURL)
}
