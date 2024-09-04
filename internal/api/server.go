package api

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/go-playground/validator/v10"
	"github.com/xavesen/search-admin/internal/config"
	"github.com/xavesen/search-admin/internal/middleware"
	"github.com/xavesen/search-admin/internal/utils"
	"github.com/xavesen/search-admin/internal/storage"
	ut "github.com/go-playground/universal-translator"
)

type Server struct {
	listenAddr	string
	storage 	storage.Storage
	config		*config.Config
	router 		*mux.Router
	validator	*validator.Validate
	translator	*ut.Translator
}

func NewServer(listenAddr string, storage storage.Storage, config *config.Config) *Server {
	log.Debug("Initializing server")

	validate, translator := utils.NewValidator()

	server := Server{
		listenAddr: listenAddr,
		storage: 	storage,
		config: 	config,
		router: 	mux.NewRouter(),
		validator:	validate,
		translator: translator,
	}
	
	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	log.Debug("Initializing routes")

	s.router.Use(middleware.Logging)

	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
	s.router.HandleFunc("/user", s.CreateUser).Methods("POST")
	s.router.HandleFunc("/users", s.GetAllUsers).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9a-z]+}", s.GetUserById).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9a-z]+}", s.DeleteUser).Methods("DELETE")
	s.router.HandleFunc("/user/{id:[0-9a-z]+}", s.UpdateUser).Methods("PUT")
	s.router.HandleFunc("/filter", s.CreateFilter).Methods("POST")
	s.router.HandleFunc("/filters", s.GetAllFilters).Methods("GET")
	s.router.HandleFunc("/filter/{id:[0-9a-z]+}", s.DeleteFilter).Methods("DELETE")
}
 
func (s *Server) Start() error {
	log.Infof("Starting listening on %s", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, s.router)
}