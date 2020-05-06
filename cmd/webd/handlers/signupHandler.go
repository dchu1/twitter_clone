package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"google.golang.org/grpc/status"
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

	_, err = UserServiceClient.CreateUser(r.Context(), &userpb.AccountInformation{FirstName: reqMessage.Firstname, LastName: reqMessage.Lastname, Email: reqMessage.Email})

	if err != nil {
		e, _ := status.FromError(err)
		if e.Message() == "duplicate email" {
			APIResponse(w, r, http.StatusConflict, "Email Id exists", make(map[string]string))
		} else {
			APIResponse(w, r, http.StatusInternalServerError, "Database not responding", make(map[string]string))
		}
		return
	}
	_, err = AuthClient.AddCredential(r.Context(), &authpb.UserCredential{Username: reqMessage.Email, Password: reqMessage.Password})
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Database not responding", make(map[string]string))
		return
	}
	APIResponse(w, r, http.StatusCreated, "Signup successful", make(map[string]string))
}
