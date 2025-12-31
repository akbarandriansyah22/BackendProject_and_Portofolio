package common

import (
	"encoding/json"
	"net/http"
)

//Respon standar untuk API
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponJSON mengirimkan response dalam format JSON
func ResponJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	response := APIResponse{
		Status:  code,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// ResponWithError mengirimkan response error dalam format JSON
func ResponWithError(w http.ResponseWriter, code int, message string) {
	ResponJSON(w, code, message, nil)
}