package api

import (
	"net/http"

	"github.com/xavesen/search-admin/internal/config"
	"github.com/xavesen/search-admin/internal/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr	string
	storage 	storage.Storage
	config		*config.Config
	router 		*mux.Router
}

func NewServer(listenAddr string, storage storage.Storage, config *config.Config) *Server {
	server := Server{
		listenAddr: listenAddr,
		storage: 	storage,
		config: 	config,
		router: 	mux.NewRouter(),
	}
	server.initialiseRoutes()
	return &server
}

func (s *Server) initialiseRoutes() {
	s.router.HandleFunc("/ping", s.Ping).Methods("GET")
	s.router.HandleFunc("/user", s.CreateUser).Methods("POST")
	s.router.HandleFunc("/users", s.GetAllUsers).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9a-z]+}", s.GetUserById).Methods("GET")
}
 
func (s *Server) Start() error {
	http.HandleFunc("/ping", s.Ping)
	return http.ListenAndServe(s.listenAddr, s.router)
}