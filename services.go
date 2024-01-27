package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserData struct {
	Username string `json:"username"`
	Age      int    `json:"age"`
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved GET reuqest from %s for %s", r.RemoteAddr, r.URL.Path)

	responseMessage := fmt.Sprintf("Hello, this is your Go server!")
	fmt.Fprint(w, responseMessage)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received POST request from %s for %s", r.RemoteAddr, r.URL.Path)
	var userData UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Data received: Username - %s, Age - %d", userData.Username, userData.Age)

	// Send a response to the client
	responseMessage := fmt.Sprintf("Data received is Username: %s and Age: %d", userData.Username, userData.Age)
	fmt.Fprintf(w, responseMessage)
}
