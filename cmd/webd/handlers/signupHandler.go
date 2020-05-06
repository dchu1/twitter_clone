package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	signuppb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/signup/signuppb"
)

// Signup is the handler for /signup. It is used for creating new users.
func Signup(w http.ResponseWriter, r *http.Request) {
	reqMessage := handlermodels.CreateUserRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Error in reading request", make(map[string]string)) // send data to client side
		return
	}
	err = json.Unmarshal(b, &reqMessage)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Error in unmarshalling", make(map[string]string)) // send data to client side
		return
	}

	// _, err = UserServiceClient.CreateUser(r.Context(), &userpb.AccountInformation{FirstName: reqMessage.Firstname, LastName: reqMessage.Lastname, Email: reqMessage.Email})
	// if err != nil {
	// 	APIResponse(w, r, http.StatusInternalServerError, "Signup unsuccessful", make(map[string]string))
	// 	return
	// }
	// _, err = AuthClient.AddCredential(r.Context(), &authpb.UserCredential{Username: reqMessage.Email, Password: reqMessage.Password})
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	_, err = SignupServiceClient.Signup(r.Context(), &signuppb.AccountInformation{FirstName: reqMessage.Firstname, LastName: reqMessage.Lastname, Email: reqMessage.Email, Password: reqMessage.Password})
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Error signing up: %v", err), make(map[string]string)) // send data to client side
		return
	}
	APIResponse(w, r, http.StatusCreated, "Signup successful", make(map[string]string))
}
