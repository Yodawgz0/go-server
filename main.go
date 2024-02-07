package main

import (
	"log"
	"net/http"
)

// Define a type that implements the http.Handler interface
type tableHandler struct{}

// Implement the ServeHTTP method for the tableHandler type
func (h *tableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check the request method and handle accordingly
	switch r.Method {
	case http.MethodGet:
		handleReadGetRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Use the http.HandleFunc function for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r)
		case http.MethodPost:
			handlePostRequest(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	th := &tableHandler{}
	http.Handle("/readTables", th)

	// Start the HTTP server
	log.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
