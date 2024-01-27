package main

import (
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
}
