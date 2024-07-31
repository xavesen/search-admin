package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode		int		`json:"statusCode"`
	ErrorMessage	string	`json:"errorMessage"`
	Data			any		`json:"data"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, errorMessage string, data any) error {
	resp := Response{
		StatusCode: statusCode,
		ErrorMessage: errorMessage,
		Data: data,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(resp)
}