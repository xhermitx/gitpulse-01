package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/service/user"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter.PathPrefix("/user").Subrouter())

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
