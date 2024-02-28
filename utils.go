package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func verifyCaptchaToken(token string) (bool, error) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret":   {"6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI"},
		"response": {token},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var result struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	fmt.Print((result.Success))
	return true, nil
}

func verifyTokenHandler(_ http.ResponseWriter, r *http.Request) (bool, error) {
	cookie, err := r.Cookie("authToken")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, err
		}
		return false, err
	}
	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, err
		}
		return false, err
	}
	if !token.Valid {
		return false, err
	}
	expirationTime := time.Unix(claims.ExpiresAt, 0)
	if expirationTime.Before(time.Now()) {
		return false, err
	}
	return true, nil
}
