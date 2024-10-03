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

	_ "github.com/xhermitx/gitpulse-01/backend/docs"
)

type message map[string]any

type Handler struct {
	userStore types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		userStore: store,
	}
}

// TODO: Implement Error Wrapper

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.HandleRegister).Methods("POST")
	router.HandleFunc("/login", h.HandleLogin).Methods("POST")
	router.HandleFunc("/update", auth.AuthMiddleware(h.HandleUpdate, h.userStore)).Methods("PATCH")
	router.HandleFunc("/delete", auth.AuthMiddleware(h.HandleDelete, h.userStore)).Methods("POST")
}

// Login
//
//	@Summary		Login
//	@Description	Login using Email/Username
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body	types.Credentials	true	"Login"	example(types.Credentials)
//	@Success		200			{json}	string				"Success"
//	@Router			/api/v1/auth/login [post]
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var credentials types.Credentials
	var user *types.User
	var err error

	if err = utils.ParseRequestBody(r, &credentials); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}

	// Support Login via Username or Email
	switch true {
	case credentials.Email != "":
		user, err = h.userStore.FindUserByEmail(credentials.Email)
	case credentials.Username != "":
		user, err = h.userStore.FindUserByUsername(credentials.Username)
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

// Register
//
//	@Summary		Register
//	@Description	Register a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		201				{json}	string		"Success"
//	@Param			user_details	body	types.User	true	"Create Account (userId not required)"	example(types.User)
//	@Router			/api/v1/auth/register [post]
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

	if err := h.userStore.CreateUser(user); err != nil {
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

// Update Account
//
//	@Summary		Update Account
//	@Description	Update Account Details
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user_details	body	types.User	true	"Update Account"	example(types.User)
//	@Success		200				{json}	string		"Success"
//	@Router			/api/v1/auth/update [patch]
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := utils.ParseRequestBody(r, &user); err != nil {
		utils.ErrResponseWriter(w, http.StatusBadRequest, err)
		return
	}
	// Return if user does not exist
	if !h.checkUserExists(w, user.UserId) {
		return
	}

	if err := h.userStore.UpdateUser(user); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}

	payload := message{
		"message":      "Updated User Successfully",
		"user_details": user,
	}
	utils.ResponseWriter(w, http.StatusOK, payload)
}

// Delete account
//
//	@Summary		Delete Account
//	@Description	Delete account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{json}	string	"Success"
//	@Router			/api/v1/auth/delete [post]
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)
	// Return if user does not exist
	if !h.checkUserExists(w, userId) {
		return
	}

	if err := h.userStore.DeleteUser(userId); err != nil {
		utils.ErrResponseWriter(w, http.StatusInternalServerError, err)
		return
	}
	payload := message{
		"message": "Deleted User Successfully",
	}
	utils.ResponseWriter(w, http.StatusOK, payload)
}

func (h *Handler) checkUserExists(w http.ResponseWriter, userId string) bool {
	res, err := h.userStore.FindUserById(userId)
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
