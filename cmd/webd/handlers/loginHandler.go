package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

// Login is the handler for /login
func Login(w http.ResponseWriter, r *http.Request) {
	var user handlermodels.LoginRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		APIResponse(w, r, http.StatusInternalServerError, "Error in reading request", make(map[string]string)) // send data to client side
		return
	}
	err = json.Unmarshal([]byte(body), &user)
	if err != nil {
		APIResponse(w, r, http.StatusUnauthorized, "Error in unmarshalling request", make(map[string]string)) // send data to client side
		return
	}
	//code to check if user exists
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := AuthClient.CheckAuthentication(ctx, &authpb.UserCredential{Username: user.Email, Password: user.Password})
	fmt.Printf("res: %v, err: %v", res, err)
	if err != nil {
		APIResponse(w, r, http.StatusUnauthorized, fmt.Sprintf("Login unsuccessful: %s", err), make(map[string]string)) // send data to client side
		return
	}

	if res.Authenticated {
		user, err := UserServiceClient.GetUserIdByUsername(r.Context(), &userpb.UserName{Email: user.Email})
		if err != nil {
			APIResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Login unsuccessful: %s", err), make(map[string]string)) // send data to client side
			return
		}

		authToken, err := AuthClient.GetAuthToken(ctx, &authpb.UserId{UserId: user.UserId})
		if err != nil {
			APIResponse(w, r, http.StatusUnauthorized, fmt.Sprintf("Login unsuccessful: %s", err), make(map[string]string)) // send data to client side
		}

		cookie := http.Cookie{Name: "sessionId", Value: url.QueryEscape(authToken.Token), Path: "/", HttpOnly: true}
		http.SetCookie(w, &cookie)

		APIResponse(w, r, http.StatusOK, "Login successful", make(map[string]string)) // send data to client side
	} else {
		APIResponse(w, r, http.StatusUnauthorized, "Login unsuccessful: User not authorized", make(map[string]string)) // send data to client side
	}

}
