package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xavesen/search-admin/internal/models"
	"github.com/xavesen/search-admin/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	users, err := s.storage.GetAllUsers(ctx)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", users)
}

func (s *Server) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "No user id provided", nil)
		return
	}

	ctx := context.TODO()
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == primitive.ErrInvalidHex {
			utils.WriteJSON(w, r, http.StatusNotFound, false, "No user with such id", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		}
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", user)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser *models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newUser) ; err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	ctx := context.TODO()
	newUser, err := s.storage.CreateUser(ctx, newUser)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusCreated, true, "", newUser)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "No user id provided", nil)
		return
	}

	ctx := context.TODO()
	err := s.storage.DeleteUser(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == primitive.ErrInvalidHex {
			utils.WriteJSON(w, r, http.StatusNotFound, false, "No user with such id", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		}
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", nil)
}