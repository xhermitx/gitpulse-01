package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ResponseWriter(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(msg)
}

func ErrResponseWriter(w http.ResponseWriter, status int, err error) {
	ResponseWriter(w, status, map[string]string{"error": err.Error()})
}

func ParseRequestBody(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
