package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

// Signup is the handler for /signup. It is used for creating new users.
func Signup(w http.ResponseWriter, r *http.Request) {
	reqMessage := handlermodels.CreateUserRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Signup unsuccessful", make(map[string]string)) // send data to client side
		return
	}
	err = json.Unmarshal(b, &reqMessage)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Signup unsuccessful", make(map[string]string)) // send data to client side
		return
	}
	_, err = application.AddUser(reqMessage.Firstname, reqMessage.Lastname, reqMessage.Email, reqMessage.Password)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Signup unsuccessful", make(map[string]string)) // send data to client side
		return
	}

	// Add user credentials to auth server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = AuthClient.AddCredential(ctx, &authpb.UserCredential{Username: reqMessage.Email, Password: reqMessage.Password})

	APIResponse(w, r, http.StatusCreated, "Signup successful", make(map[string]string))
}
