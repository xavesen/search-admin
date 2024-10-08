package api

import (
	"net/http"

	"github.com/xavesen/search-admin/internal/utils"
)

type PingResponse struct {
	Pong	string	`json:"pong"`
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, r, http.StatusOK, true, "", PingResponse{Pong: "pong"})
}