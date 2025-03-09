// This service generates short URLs and stores them in PostgreSQL.

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type URLShortener struct {
	db *sql.DB
}

type ShortURL struct {
	ID        int64
	ShortCode string
	LongURL   string
	CreatedAt time.Time
}

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", os.Getenv("LOCAL_DATABASE_URL") + "?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create shortener service
	shortener := &URLShortener{db: db}

	// Initialize router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/shorten", shortener.ShortenHandler).Methods("POST")

	// Start server
	log.Println("Starting URL shortener service on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s *URLShortener) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a short code
	// Store the short code in the database
	// Return the short code to the client

	// Get the long URL from the request body
	var requestBody struct {
		LongURL string `json:"long_url"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate a short code
	shortCode := generateShortCode()

	// Store the short code in the database
	var id int64
	err = s.db.QueryRow("INSERT INTO urls (short_code, long_url) VALUES ($1, $2) RETURNING id", shortCode, requestBody.LongURL).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to store short code", http.StatusInternalServerError)
		log.Printf("Failed to store short code: %v", err)
		return
	}

	log.Printf("Short code created in DB: Long URL: %s, Short Code: %s, ID: %d", requestBody.LongURL, shortCode, id)
	// Return the short code to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortCode))
}

func generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
