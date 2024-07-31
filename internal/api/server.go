package api

import (
	"net/http"

	"github.com/xavesen/search-admin/internal/config"
	"github.com/xavesen/search-admin/internal/storage"
)

type Server struct {
	listenAddr	string
	storage 	storage.Storage
	config		*config.Config
}

func NewServer(listenAddr string, storage storage.Storage, config *config.Config) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage: 	storage,
		config: config,
	}
}
 
func (s *Server) Start() error {
	http.HandleFunc("/ping", s.Ping)
	return http.ListenAndServe(s.listenAddr, nil)
}