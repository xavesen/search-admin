package utils

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextKeyReqId ContextKey = "requestId"

type Response struct {
	Success			bool	`json:"success"`
	ErrorMessage	string	`json:"errorMessage"`
	Data			any		`json:"data"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, statusCode int, success bool, errorMessage string, data any) error {
	log.WithFields(log.Fields{
		"request_id": r.Context().Value(ContextKeyReqId).(string),
		"status_code": statusCode,
		"error_message": errorMessage,
	}).Info("Responding to request")
	
	resp := Response{
		Success: success,
		ErrorMessage: errorMessage,
		Data: data,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(resp)
}