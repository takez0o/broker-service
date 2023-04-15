package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var request_payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &request_payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadGateway)
		return
	}

	user, err := app.Models.User.GetByEmail(request_payload.Email)
	fmt.Println("Payload:", request_payload.Email, request_payload.Password)
	if err != nil {
		app.errorJSON(w, errors.New("1Invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(request_payload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("2Invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
