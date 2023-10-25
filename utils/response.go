package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents the structure of the HTTP response.
type Response struct {
	Message string `json:"message"`
}

// SendResponse sends an HTTP response with the given status code and message.
func SendResponse(w http.ResponseWriter, statusCode int, message string) {
	response := Response{Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}