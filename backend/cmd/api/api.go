package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xhermitx/gitpulse-01/backend/service/job"
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
	subrouter := router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter.PathPrefix("/user").Subrouter())

	jobStore := job.NewStore(s.db)
	jobHandler := job.NewHandler(jobStore, userStore)
	jobHandler.RegisterRoutes(subrouter.PathPrefix("/job").Subrouter())

	subrouter.HandleFunc("/", Greetings)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

func Greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Greetings from gitpulse!")
}
