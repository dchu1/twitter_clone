package handlers

import (
	"encoding/json"
	"net/http"
)

// ResponseMessage is the struct for the APIResponse responses
type ResponseMessage struct {
	Status  int
	Message string
	Body    map[string]string
}

// APIResponse the function for writing a generic API response
func APIResponse(w http.ResponseWriter, r *http.Request, status int, message string, body map[string]string) {
	response := ResponseMessage{
		Status:  status,
		Message: message,
		Body:    body,
	}

	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")

	w.Write(res)
}
