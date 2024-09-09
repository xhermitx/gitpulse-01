package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/resume-parser/types"
)

const (
	// TODO: Add other drive link patterns
	GOOGLE_PATTERN = `https://drive\.google\.com/drive/folders/([0-9A-Za-z-_]+)`

	JOB_STATUS_QUEUE = "JOB_STATUS_QUEUE"
)

type APIServer struct {
	addr    string
	storage types.Drive
	queue   types.Queue
	cache   types.KVStore
}

func NewAPIServer(addr string, d types.Drive, q types.Queue) *APIServer {
	return &APIServer{
		addr:    addr,
		storage: d,
		queue:   q,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	subrouter.HandleFunc("/trigger/{provider}", s.TriggerHandler).Methods(http.MethodPost)

	log.Printf("Listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) TriggerHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	provider := vars["provider"]

	var trigger types.TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&trigger); err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	folderId, err := extractFolderID(provider, trigger.DriveLink)
	if err != nil {
		errResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	resumeList, err := s.storage.GetFileList(folderId)
	if err != nil {
		errResponseWriter(w, http.StatusInternalServerError, errors.New("couldn't fetch resumes"))
		return
	}

	// Flush and close the HTTP connection the before proceeding
	closeConnection(w)

	var wg sync.WaitGroup
	wg.Add(len(resumeList))

	for _, resumeId := range resumeList {
		go func() {
			defer wg.Done()
			s.handleResume(resumeId, trigger)
		}()
	}
	wg.Wait()
}

func (s *APIServer) handleResume(rId string, trigger types.TriggerRequest) {

	fileName, err := s.cache.Get(context.Background(), rId)
	if err != nil {
		log.Println(err)
	}

	type resumeStatus struct {
		fileName string
		status   bool
	}

	content, err := s.storage.GetFileContent(rId)
	if err != nil {
		// FIXME: Get the filename from Redis using FileId
		// 		  Update the status of the file as unparsed
		_ = s.cache.Set(context.Background(), trigger.JobId, resumeStatus{
			fileName: fileName.(string),
			status:   false,
		})
		_ = err
	}

	githubIds, err := s.storage.GetUsername(content)
	if err != nil {
		// FIXME: Get the filename from Redis using FileId
		// 		  Update the status of the file as unparsed
		_ = s.cache.Set(context.Background(), fileName, resumeStatus(resumeStatus{
			fileName: fileName.(string),
			status:   false,
		}))
		_ = err
	}

	log.Printf("\n%s : \n%s", trigger.JobId, githubIds)

	// TODO: Push the username along with JobId to RabbitMQ
	if err = s.queue.Publish(JOB_STATUS_QUEUE, types.StatusQueue{
		JobId:     trigger.JobId,
		FileId:    rId,
		GithubIDs: githubIds,
	}); err != nil {
		// FIXME: Get the filename from Redis using FileId
		// 		  Update the status of the file as unparsed
		_ = s.cache.Set(context.Background(), fileName, resumeStatus(resumeStatus{
			fileName: fileName.(string),
			status:   false,
		}))
		_ = err
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

func extractFolderID(provider, link string) (string, error) {
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

func closeConnection(w http.ResponseWriter) {
	responseWriter(w, http.StatusOK, map[string]string{
		"Message": "Successfully triggered",
	})
	w.(http.Flusher).Flush()

	// Hijack the connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		errResponseWriter(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		errResponseWriter(w, http.StatusInternalServerError, errors.New("internal error"))
		return
	}
	conn.Close()
}
