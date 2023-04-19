package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var request_payload JSONPayload
	_ = app.readJSON(w, r, &request_payload)

	event := data.LogEntry{
		Name: request_payload.Name,
		Data: request_payload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Log entry created",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
