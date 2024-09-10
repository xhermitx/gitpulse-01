package job

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/service/auth"
	"github.com/xhermitx/gitpulse-01/backend/types"
	"github.com/xhermitx/gitpulse-01/backend/utils"
)

type message map[string]any

type Handler struct {
	jobStore  types.JobStore
	userStore types.UserStore
}

func NewHandler(jobStore types.JobStore, userStore types.UserStore) *Handler {
	return &Handler{
		jobStore:  jobStore,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create", auth.AuthMiddleware(h.CreateJobHandler, h.userStore)).Methods("POST")
	router.HandleFunc("/udpate", auth.AuthMiddleware(h.UpdateJobHandler, h.userStore)).Methods("PATCH")
	router.HandleFunc("/delete", auth.AuthMiddleware(h.DeleteJobHandler, h.userStore)).Methods("DELETE")
	router.HandleFunc("/list", auth.AuthMiddleware(h.ListJobHandler, h.userStore)).Methods("GET")
	router.HandleFunc("/trigger", auth.AuthMiddleware(h.TriggerJobHandler, h.userStore)).Methods("POST")
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

	if !h.checkJobExists(w, job.JobId) {
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

	log.Println("job_id: ", job.JobId)

	if !h.checkJobExists(w, job.JobId) {
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
	// TODO: To be implemented
	// 		 Remember to provide the {cloudprovide in the request param}
}

func (h *Handler) checkJobExists(w http.ResponseWriter, jobId string) bool {
	res, err := h.jobStore.FindJobById(jobId)
	if res == nil {
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job does not exist"))
		return false
	}
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return false
	}
	return true
}
