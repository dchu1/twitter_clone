package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	reqMessage := handlermodels.CreateUserRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b, &reqMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = application.AddUser(reqMessage.Firstname, reqMessage.Lastname, reqMessage.Email, reqMessage.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	APIResponse(w, r, http.StatusCreated, "Signup successful", make(map[string]string))
}
