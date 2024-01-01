package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteStatus(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Failed to encode API response into JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func WriteError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := errResponse{
		Status:  statusCode,
		Message: errMsg,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response into JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
