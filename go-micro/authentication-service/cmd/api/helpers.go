// TODO: remove duplicate by lib package
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

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != nil && err != io.EOF {
		if err == io.ErrUnexpectedEOF {
			return errors.New("JSON data is too large and exceeds the maximum buffer size (1MB)")
		}
		return errors.New("request body contains more than one JSON object, which is not allowed")
	}
	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, header := range headers {
		for k, v := range header {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	code := http.StatusBadRequest
	if len(status) > 0 {
		code = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return app.writeJSON(w, code, payload)
}
