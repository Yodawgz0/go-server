package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r)
		case http.MethodPost:
			handlePostRequest(w, r)
		default:
			http.Error(w, "Methodd not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/readTables", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleReadGetRequest(w, r)
		default:
			http.Error(w, "Methodd not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Start the HTTP server
	log.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
