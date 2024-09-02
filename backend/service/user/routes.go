package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/types"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/delete", h.HandleDeleteUser).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	// Check if user Exists
	// If not, create a new user
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

}
