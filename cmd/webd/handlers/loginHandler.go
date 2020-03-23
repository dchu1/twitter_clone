package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("Hi")
	session.Sess = session.GlobalSessions.SessionStart(w, r)
	json.Unmarshal([]byte(body), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//code to check if user exists
	fmt.Println("Hello")
	session.Sess.Set("username", user.Email)
	session.Sess.Set("authenticated", true)
	APIResponse(w, r, 200, "Login successful", make(map[string]string)) // send data to client side
}
