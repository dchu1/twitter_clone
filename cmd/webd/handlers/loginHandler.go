package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user handlermodels.LoginRequest
	body, err := ioutil.ReadAll(r.Body)
	sess := session.GlobalSessions.SessionStart(w, r)
	json.Unmarshal([]byte(body), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//code to check if user exists
	if a.ValidateCredentials(user.Email, user.Password) {
		sess.Set("username", user.Email)
		sess.Set("authenticated", true)
		APIResponse(w, r, http.StatusOK, "Login successful", make(map[string]string)) // send data to client side
	} else {
		APIResponse(w, r, http.StatusUnauthorized, "Login unsuccessful", make(map[string]string)) // send data to client side
	}

}
