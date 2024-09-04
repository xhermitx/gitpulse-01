package user

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
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

// TODO: Implement Authentication Middleware

// TODO: Implement Error Wrapper

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/delete", h.HandleDeleteUser).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials types.Credentials
	if err := utils.ParseRequestBody(r, &credentials); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.LoginUser(credentials)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	token := utils.GenerateToken(user.UserId)
	payload := message{
		"message":      "Login Successful!",
		"user_details": user,
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Secure:   false,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	utils.ResponseWriter(w, http.StatusOK, payload)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON Payload
	var user types.User
	if err := utils.ParseRequestBody(r, &user); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	user.UserId = uuid.NewString()

	ok, err := h.checkUserExists(user.Email, user.Username)
	if ok {
		utils.ErrResponseWriter(w, http.StatusConflict, errors.New("user already exists"))
		return
	}
	if err != nil {
		utils.ResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.CreateUser(user); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message":      "User Created Successfully",
		"user_details": user,
	}

	utils.ResponseWriter(w, http.StatusCreated, payload)
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) checkUserExists(email string, username string) (bool, error) {
	res, err := h.store.FindUserByUsername(username)
	if res != nil {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	res, err = h.store.FindUserById(email)
	if res != nil {
		return true, nil
	}

	return false, err
}
