package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("cassandra")

type Credentials struct {
	Email        string `json:"email"`
	CaptchaToken string `json:"captchaToken"`
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	captchaToken := creds.CaptchaToken
	if captchaToken == "" {
		http.Error(w, "reCAPTCHA token is required", http.StatusBadRequest)
		return
	}

	// Verify reCAPTCHA token
	isValid, err := verifyCaptchaToken(captchaToken)
	if err != nil {
		http.Error(w, "Failed to verify reCAPTCHA token", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Failed to verify reCAPTCHA token", http.StatusUnauthorized)
		return
	}
	email := creds.Email
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	fmt.Print(expirationTime)

	http.SetCookie(w, &http.Cookie{
		Name:    "authToken",
		Value:   tokenString,
		Expires: expirationTime,
	})

	fmt.Fprint(w, "User Successfuly Logged In!")
}
func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "authToken",
		Value: "",
		// Path:     "/", // specify the path if necessary
		// Domain:   "example.com", // specify the domain if necessary
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		// Secure:   true, // set to true if your application is served over HTTPS
	})

	fmt.Fprint(w, "User logged out successfully")
}
