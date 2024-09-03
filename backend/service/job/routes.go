package job

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
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

}

func (h *Handler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteJobHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListJobHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) TriggerJobHandler(w http.ResponseWriter, r *http.Request) {

}
