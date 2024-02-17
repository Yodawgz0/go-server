package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
