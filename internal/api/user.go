package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/xavesen/search-admin/internal/models"
	"github.com/xavesen/search-admin/internal/utils"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser *models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser) ; err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	ctx := context.TODO()
	newUser, err := s.storage.CreateUser(ctx, newUser)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, true, "", newUser)
}