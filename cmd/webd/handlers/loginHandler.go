package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

// Login is the handler for /login
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := AuthClient.CheckAuthentication(ctx, &authpb.UserCredential{Username: user.Email, Password: user.Password})
	fmt.Println(AuthClient)

	if err != nil {
		APIResponse(w, r, http.StatusUnauthorized, "Login unsuccessful", make(map[string]string)) // send data to client side
	}
	fmt.Println(res, err)
	if res.Authenticated {
		user, err := application.GetUserByUsername(user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sess.Set("userId", user.Id)
		sess.Set("username", user.Email)
		sess.Set("authenticated", true)
		APIResponse(w, r, http.StatusOK, "Login successful", make(map[string]string)) // send data to client side
	} else {
		APIResponse(w, r, http.StatusUnauthorized, "Login unsuccessful", make(map[string]string)) // send data to client side
	}

}
