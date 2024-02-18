package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", middleware(handleRequest))
	http.HandleFunc("/userLogin", middleware(userLoginHandler))
	http.HandleFunc("/userLogout", middleware(userLogoutHandler))

	log.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
