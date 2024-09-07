package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/resume-parser/drive"
	"github.com/xhermitx/gitpulse-01/resume-parser/types"
)

const (
	google_pattern = `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`
	// TODO: Add other drive link patterns
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

	subrouter.HandleFunc("/trigger/{provider}", s.TriggerHandler).Methods(http.MethodPost)

	return http.ListenAndServe(s.addr, subrouter)
}

func (s *APIServer) TriggerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	provider := vars["provider"]

	var payload types.TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	folderId, err := extractFolderID(provider, payload.DriveLink)
	if err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	var driver types.Drive

	switch provider {
	case "google":
		driver, err = drive.NewGoogleDrive(folderId)
		if err != nil {
			errResponseWriter(w, http.StatusInternalServerError, err)
			return
		}
		// TODO: Add other Cloud Providers' logic
	}

	resumeList, err := driver.GetFileList()
	if err != nil {
		errResponseWriter(w, http.StatusInternalServerError, errors.New("couldn't fetch resumes"))
		return
	}

	w.(http.Flusher).Flush()

	var (
		wg sync.WaitGroup
	)

	wg.Add(len(resumeList))
	for _, resumeId := range resumeList {
		go func() {
			defer wg.Done()
			content, err := driver.GetFileContent(resumeId)
			if err != nil {
				// FIXME: Get the filename from Redis using FileId
				// 		  Update the status of the file as unparsed
				_ = err
			}

			githubIds, err := driver.GetUsername(content)
			if err != nil {
				// FIXME: Get the filename from Redis using FileId
				// 		  Update the status of the file as unparsed
				_ = err
			}

			log.Printf("\n%s : \n%s", payload.JobId, githubIds)

			// TODO: Push the username along with JobId to RabbitMQ
		}()
	}

	wg.Wait()
}

func responseWriter(w http.ResponseWriter, status int, msg any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(msg)
}

func errResponseWriter(w http.ResponseWriter, status int, err error) {
	responseWriter(w, status, map[string]string{"error": err.Error()})
}

func extractFolderID(provider, link string) (string, error) {
	var re *regexp.Regexp

	switch provider {
	case "google":
		re = regexp.MustCompile(google_pattern)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			// First part is the entire match, second is the captured group
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("folder ID not found in link")
}
