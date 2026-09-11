package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}
