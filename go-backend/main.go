package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

// PostgreSQL connection details
const (
	host     = "localhost"
	port     = 5432
	user     = "myuser1"
	password = "mypassword"
	dbname   = "database2"
)

var db *sql.DB

func main() {
	// Database connection setup
	var err error
	psqlInfo := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	psqlInfo = fmt.Sprintf(psqlInfo, host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},           // Allow all origins
		AllowedMethods: []string{"GET", "POST"}, // Allowed methods
		AllowedHeaders: []string{"Content-Type"},
	})

	// HTTP routes
	http.HandleFunc("/post-endpoint", handlePost)
	http.HandleFunc("/get-endpoint", handleGet)

	// CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	// Start server
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Logic to store data in PostgreSQL
	_, err = db.Exec("INSERT INTO user_data (data) VALUES ($1)", data["data"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data) // Sending back the data as response
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the URL query
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Error parsing query: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := queryParams.Get("data")
	if data == "" {
		http.Error(w, "Missing data parameter", http.StatusBadRequest)
		return
	}

	var responseMessage string
	err = db.QueryRow("SELECT data FROM user_data WHERE data = $1", data).Scan(&responseMessage)
	if err != nil {
		// Handle the error here
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, responseMessage)
}
