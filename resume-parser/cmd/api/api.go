package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/resume-parser/types"
	"github.com/xhermitx/gitpulse-01/resume-parser/utils"
)

const (
	QUEUE__JOB_STATUS    = "JOB_STATUS_QUEUE"
	CACHE__JOB_STATUS    = "JOB_STATUS_CACHE"
	CACHE__FAILED_RESUME = "FAILED_RESUME_CACHE"
)

type APIServer struct {
	addr    string
	storage types.Drive
	queue   types.Queue
	cache   types.KVStore
}

func NewAPIServer(addr string, d types.Drive, q types.Queue, c types.KVStore) *APIServer {
	return &APIServer{
		addr:    addr,
		storage: d,
		queue:   q,
		cache:   c,
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
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	// Extract folderId from given drive link based on the cloud provider
	folderId, err := utils.ExtractFolderID(provider, trigger.DriveLink)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	// Get a map containing fileId: fileName
	resumeList, err := s.storage.GetFileList(folderId)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, errors.New("couldn't fetch files"))
		return
	}

	// Flush and close the HTTP connection the before proceeding
	utils.CloseConnection(w)

	var wg sync.WaitGroup
	wg.Add(len(resumeList))

	for fileId, fileName := range resumeList {
		go func() {
			defer wg.Done()
			s.handleResume(fileId, fileName, trigger)
		}()
	}
	wg.Wait()

	// Update the Job Queue status as True
	if err = s.queue.Publish(QUEUE__JOB_STATUS, types.JobQueue{
		JobId:  trigger.JobId,
		Status: true,
	}); err != nil {
		// FIXME: Handle this better
		log.Println(err)
	}
}

func (s *APIServer) handleResume(fId, fName string, trigger types.TriggerRequest) {
	content, err := s.storage.GetFileContent(fId)
	if err != nil {
		if err := s.cacheFile(fName, trigger.JobId); err != nil {
			log.Printf("\nfailed parsing %s : %v", fName, err)
		}
		log.Println(err)
	}

	githubIds, err := s.storage.GetUsername(content)
	if err != nil {
		if err := s.cacheFile(fName, trigger.JobId); err != nil {
			log.Printf("\nfailed parsing %s : %v", fName, err)
		}
		log.Println(err)
	}

	log.Printf("\n%s : \n%s", trigger.JobId, githubIds)

	// Push the username along with JobId to RabbitMQ
	if err = s.queue.Publish(QUEUE__JOB_STATUS, types.JobQueue{
		JobId:     trigger.JobId,
		Filename:  fName,
		GithubIDs: githubIds,
		Status:    false,
	}); err != nil {
		if err := s.cacheFile(fName, trigger.JobId); err != nil {
			log.Printf("\nfailed parsing %s : %v", fName, err)
		}
		log.Println(err)
	}
}

func (s *APIServer) cacheFile(filename string, jobId string) error {
	var (
		key   = CACHE__FAILED_RESUME + jobId
		files = []string{}
	)

	res, err := s.cache.Get(context.Background(), jobId)
	if err != nil {
		return err
	}
	if res != "" {
		err := json.Unmarshal([]byte(res), &files)
		if err != nil {
			return err
		}
	}

	files = append(files, filename)
	if err := s.cache.Set(context.Background(), key, files, 0); err != nil {
		return err
	}
	return nil
}
