package job

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
	"github.com/xhermitx/gitpulse-01/backend/utils"
)

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

	res, err := h.store.FindJobByName(job.Name)
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

	if err = utils.ResponseWriter(w, http.StatusCreated, map[string]string{"message": "Created Job Successfully"}); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	utils.ErrResponseWriter(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func (h *Handler) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {
	var job types.Job
	if err := utils.ParseRequestBody(r, &job); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.store.FindJobByName(job.Name)
	switch true {
	case res != nil:
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("job already exists"))
		return
	case err != nil:
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) ListJobHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) TriggerJobHandler(w http.ResponseWriter, r *http.Request) {

}
