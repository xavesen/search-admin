package api

import (
	"net/http"

	"github.com/xavesen/search-admin/internal/utils"
)

type PingResponse struct {
	Pong	string	`json:"ping"`
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, true, "", PingResponse{Pong: "pong"})
}