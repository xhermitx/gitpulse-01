package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/config"
	"github.com/xhermitx/gitpulse-01/backend/service/auth"
	"github.com/xhermitx/gitpulse-01/backend/types"
	"github.com/xhermitx/gitpulse-01/backend/utils"
	results "github.com/xhermitx/gitpulse-results"
)

type message map[string]any

type Handler struct {
	jobStore       types.JobStore
	userStore      types.UserStore
	candidateStore types.CandidateStore
}

func NewHandler(jobStore types.JobStore, userStore types.UserStore, candidateStore types.CandidateStore) *Handler {
	return &Handler{
		jobStore:       jobStore,
		userStore:      userStore,
		candidateStore: candidateStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create", auth.AuthMiddleware(h.CreateJobHandler, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/udpate", auth.AuthMiddleware(h.UpdateJobHandler, h.userStore)).Methods(http.MethodPatch)
	router.HandleFunc("/delete", auth.AuthMiddleware(h.DeleteJobHandler, h.userStore)).Methods(http.MethodDelete)
	router.HandleFunc("/list", auth.AuthMiddleware(h.ListJobHandler, h.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/trigger", auth.AuthMiddleware(h.TriggerJobHandler, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/result/{count}", auth.AuthMiddleware(h.ResultHandler, h.userStore)).Methods(http.MethodGet)
}

func (h *Handler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.Job
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	// Generate a new ID since it is not received in the request
	job.JobId = uuid.NewString()
	job.UserId = r.Context().Value(types.UserContext("user_id")).(string)

	if err := h.jobStore.CreateJob(job); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message":     "Job Created Successfully",
		"job_details": job,
	}
	if err := utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.Job
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	if _, ok := h.checkJobExists(w, job.JobId); !ok {
		return
	}

	if err := h.jobStore.UpdateJob(job); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Job Updated Successfully",
	}
	if err := utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.DeleteJobPayload
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	if _, ok := h.checkJobExists(w, job.JobId); !ok {
		return
	}

	if err := h.jobStore.DeleteJob(job.JobId); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Job Deleted Successfully",
	}
	if err := utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) ListJobHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(types.UserContext("user_id")).(string)

	jobs, err := h.jobStore.ListJobs(id)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Fetched Job List",
		"jobs":    jobs,
	}
	utils.ResponseWriter(w, http.StatusFound, payload)
}

func (h *Handler) TriggerJobHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.TriggerJobPayload
	if err := utils.ParseRequestBody(r, &payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	job, ok := h.checkJobExists(w, payload.JobId)
	if !ok {
		return
	}

	payload.DriveLink = job.DriveLink
	body, err := json.Marshal(payload)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, config.Envs.ParserURL, bytes.NewBuffer(body))
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("Status: ", res.StatusCode)
		utils.ErrResponseWriter(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	defer res.Body.Close()

	utils.ResponseWriter(w, http.StatusOK, message{
		"message": "Job Profiling Started!",
	})
}

func (h *Handler) ResultHandler(w http.ResponseWriter, r *http.Request) {
	var payload types.JobResultPayload
	if err := utils.ParseRequestBody(r, &payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	if _, ok := h.checkJobExists(w, payload.JobId); !ok {
		return
	}

	candidateList, err := h.candidateStore.GetCandidateList(payload.JobId)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}

	topCanidates := results.TopNCandidates(candidateList, payload.Count)

	utils.ResponseWriter(w, http.StatusOK, message{
		"message": "Result Candidates",
		"List":    topCanidates,
	})
}

func (h *Handler) checkJobExists(w http.ResponseWriter, jobId string) (*types.Job, bool) {
	res, err := h.jobStore.FindJobById(jobId)
	if res == nil {
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job does not exist"))
		return nil, false
	}
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return nil, false
	}
	return res, true
}
