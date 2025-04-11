package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse defines the structure of standard API responses
type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SendError sends a structured error JSON response
func SendError(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  false,
		Message: message,
		Error:   err,
	})
}

// SendSuccess sends a structured success JSON response
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
		Error:   nil,
	})
}
