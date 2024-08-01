package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success			bool	`json:"success"`
	ErrorMessage	string	`json:"errorMessage"`
	Data			any		`json:"data"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, success bool, errorMessage string, data any) error {
	resp := Response{
		Success: success,
		ErrorMessage: errorMessage,
		Data: data,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(resp)
}