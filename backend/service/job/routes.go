package job

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
	"github.com/xhermitx/gitpulse-01/backend/utils"
)

type message map[string]any

type Handler struct {
	store types.JobStore
}

func NewHandler(store types.JobStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create", h.CreateJobHandler).Methods("POST")
	router.HandleFunc("/udpate", h.UpdateJobHandler).Methods("PATCH")
	router.HandleFunc("/delete", h.DeleteJobHandler).Methods("DELETE")
	router.HandleFunc("/list", h.ListJobHandler).Methods("GET")
	router.HandleFunc("/trigger", h.TriggerJobHandler).Methods("POST")
}

func (h *Handler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.Job
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	// Generate a new ID since it is not received in the request
	job.JobId = uuid.NewString()

	res, err := h.store.FindJobById(job.JobId)
	switch true {
	case res != nil:
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job already exists"))
		return
	case err != nil:
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	if err = h.store.CreateJob(job); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message":     "Job Created Successfully",
		"job_details": job,
	}

	if err = utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {

	var job types.Job
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.store.FindJobById(job.JobId)
	switch true {
	case res == nil:
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job does not exist"))
		return
	case err != nil:
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	if err = h.store.UpdateJob(job); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Job Updated Successfully",
	}

	if err = utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.DeleteJobPayload
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.store.FindJobById(job.JobId)
	switch true {
	case res != nil:
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job already exists"))
		return
	case err != nil:
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	if err = h.store.DeleteJob(job.JobId); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Job Deleted Successfully",
	}

	if err = utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) ListJobHandler(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("user_id").(string)

	jobs, err := h.store.ListJobs(id)
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
}
