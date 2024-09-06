package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/resume-parser/drive/gdrive"
	"github.com/xhermitx/gitpulse-01/resume-parser/types"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	subrouter.HandleFunc("/trigger/{driver}", s.TriggerHandler).Methods(http.MethodPost)

	return http.ListenAndServe(s.addr, subrouter)
}

func (s *APIServer) TriggerHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.TriggerRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	folderId, err := extractFolderID(payload.DriveLink)
	if err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	gdrive := gdrive.NewGoogleDrive(folderId)

	resumeList, err := gdrive.GetFileList()
	if err != nil {
		errResponseWriter(w, http.StatusInternalServerError, errors.New("couldn't fetch resumes"))
		return
	}

	w.(http.Flusher).Flush()

	for _, resume := range resumeList {
		go func() {
			content, err := gdrive.GetFileContent(resume)
			if err != nil {
				// TODO: Add the filename to Redis
				_ = err
			}

			username, err := gdrive.GetUsername(bytes.NewBuffer(content))
			_ = username

			// TODO: Push the username along with JobId to RabbitMQ
		}()
	}
}

func responseWriter(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(msg)
}

func errResponseWriter(w http.ResponseWriter, status int, err error) {
	responseWriter(w, status, map[string]string{"error": err.Error()})
}

func extractFolderID(link string) (string, error) {

	pattern := `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(link)

	if len(matches) > 1 {
		// THE FIRST MATCH IS THE ENTIRE MATCH, AND THE SECOND IS THE CAPTURED GROUP
		return matches[1], nil
	}

	return "", fmt.Errorf("folder ID not found in link")
}
