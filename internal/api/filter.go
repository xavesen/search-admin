package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/go-playground/validator/v10"
	"github.com/xavesen/search-admin/internal/models"
	"github.com/xavesen/search-admin/internal/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
)

func (s *Server) CreateFilter(w http.ResponseWriter, r *http.Request) {
	var newFilter *models.Filter

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newFilter) ; err != nil {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Invalid request payload", nil)
		return
	}

	err := s.validator.Struct(newFilter)
	if err != nil {
		logErrorString := "User input validation error:"
		errorString := "Bad request:"
		for i, err := range err.(validator.ValidationErrors) {
			if i != 0 {
				errorString = errorString + ","
				logErrorString = logErrorString + ";"
			}
			errorString = errorString + " " + err.Translate(*s.translator)
			logErrorString = logErrorString + " " + err.Error()
		}
		log.WithFields(log.Fields{
			"request_id": r.Context().Value(utils.ContextKeyReqId),
			"method": r.Method,
			"url_path": r.URL.Path,
		}).Warning(logErrorString)
		utils.WriteJSON(w, r, http.StatusBadRequest, false, errorString, nil)
		return
	}

	_, err = regexp.Compile(newFilter.Regex)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": r.Context().Value(utils.ContextKeyReqId),
			"method": r.Method,
			"url_path": r.URL.Path,
		}).Warningf("Error parsing regular expression '%s' passed by user: %s", newFilter.Regex, err)
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "Bad request: regex must be a regular expression accepted by RE2", nil)
		return
	}

	ctx := context.TODO()
	newFilter, err = s.storage.CreateFilter(ctx, newFilter)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusCreated, true, "", newFilter)
}

func (s *Server) GetAllFilters(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	filters, err := s.storage.GetAllFilters(ctx)
	if err != nil {
		utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", filters)
}

func (s *Server) DeleteFilter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteJSON(w, r, http.StatusBadRequest, false, "No filter id provided", nil)
		return
	}

	ctx := context.TODO()
	err := s.storage.DeleteFilter(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == primitive.ErrInvalidHex {
			utils.WriteJSON(w, r, http.StatusNotFound, false, "No filter with such id", nil)
		} else {
			utils.WriteJSON(w, r, http.StatusInternalServerError, false, "Internal server error", nil)
		}
		return
	}

	utils.WriteJSON(w, r, http.StatusOK, true, "", nil)
}