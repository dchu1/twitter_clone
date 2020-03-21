package handlers

import (
	"encoding/json"
	"net/http"
)

type ResponseMessage struct {
	Status  int
	Message string
	Body    map[string]string
}

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
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Endcoding, Content-Type, Content-Length, Authorization, X-CSRF-token")
	// w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
	// fmt.Println(w)
	w.Write(res)
}
