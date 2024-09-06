package user

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
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

// TODO: Implement Error Wrapper

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/update", h.HandleUpdate).Methods("PATCH")
	router.HandleFunc("/delete", h.HandleDelete).Methods("POST")
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials types.Credentials
	if err := utils.ParseRequestBody(r, &credentials); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	var user *types.User
	var err error

	// Support Login via Username or Email
	switch true {
	case credentials.Email != "":
		user, err = h.store.FindUserByEmail(credentials.Email)
	case credentials.Username != "":
		user, err = h.store.FindUserByUsername(credentials.Username)
	}
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	log.Println("user found", user.Username)

	if !auth.ComparePassword([]byte(user.Password), []byte(credentials.Password)) {
		log.Println("passwords do not match")
		utils.ErrResponseWriter(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	token, err := auth.GenerateToken(user.UserId)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}
	payload := message{
		"message":      "Login Successful!",
		"user_details": user,
	}

	// TODO: Enable HTTPS
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
	var user types.User
	if err := utils.ParseRequestBody(r, &user); err != nil {
		log.Println("Error while Parsing request")
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	// Generate hash using the plain password in the request Body
	hashed, err := auth.HashedPassword(user.Password)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}

	// Assign a new uuid and generate a hashed password
	// Since it is not accepted from the request body
	user.UserId = uuid.NewString()
	user.Password = string(hashed)

	// FIXME:
	// Check if username, email already exists.
	// Although currently it will be handled on the frontend

	if err := h.store.CreateUser(user); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	log.Println("user created")

	payload := message{
		"message":      "User Created Successfully",
		"user_details": user,
	}

	if err := utils.ResponseWriter(w, http.StatusCreated, payload); err != nil {
		log.Println("Error sending response")
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := utils.ParseRequestBody(r, &user); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	if !h.checkUserExists(w, user.UserId) {
		// Return if user does not exist
		return
	}

	if err := h.store.UpdateUser(user); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message":      "Updated User Successfully",
		"user_details": user,
	}

	utils.ResponseWriter(w, http.StatusOK, payload)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var user types.DeleteUserPayload
	if err := utils.ParseRequestBody(r, &user); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	if !h.checkUserExists(w, user.UserId) {
		// Return if user does not exist
		return
	}

	if err := h.store.DeleteUser(user.UserId); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message": "Deleted User Successfully",
	}
	utils.ResponseWriter(w, http.StatusOK, payload)
}

func (h *Handler) checkUserExists(w http.ResponseWriter, userId string) bool {
	res, err := h.store.FindUserById(userId)
	if err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return false
	}
	if res == nil {
		utils.ErrResponseWriter(w, http.StatusNotFound, errors.New("user not found"))
		return false
	}

	return true
}
