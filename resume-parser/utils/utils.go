package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

const (
	GOOGLE_PATTERN = `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`
	// TODO: Add other drive link patterns
)

func ResponseWriter(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(msg)
}

func ErrResponseWriter(w http.ResponseWriter, status int, err error) {
	ResponseWriter(w, status, map[string]string{"error": err.Error()})
}

func ExtractFolderID(provider, link string) (string, error) {
	var re *regexp.Regexp

	switch provider {
	case "google":
		re = regexp.MustCompile(GOOGLE_PATTERN)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			// First part is the entire match, second is the captured group
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("folder ID not found in link")
}

func CloseConnection(w http.ResponseWriter) {
	ResponseWriter(w, http.StatusOK, map[string]string{
		"Message": "Successfully triggered",
	})
	w.(http.Flusher).Flush()

	// Hijack the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		ErrResponseWriter(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		ErrResponseWriter(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}
	conn.Close()
}
