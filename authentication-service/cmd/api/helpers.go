package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	max_bytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if headers != nil {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	status_code := http.StatusBadGateway

	if len(status) > 0 {
		status_code = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, status_code, payload)
}
