package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

// FollowCreateHandler is the Handler for following a user
func FollowCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.FollowRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Get userId from the session cookie
		sess := session.GlobalSessions.SessionQuery(w, r)
		application.FollowUser(sess.Get("userId").(uint64), reqMessage.UserId)
		APIResponse(w, r, 200, "User followed", make(map[string]string)) // send data to client side
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

// FollowDestroyHandler is the Handler for unfollowing a user
func FollowDestroyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.FollowRequest{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(body, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Get userId from the session cookie
		sess := session.GlobalSessions.SessionQuery(w, r)
		application.UnFollowUser(sess.Get("userId").(uint64), reqMessage.UserId)
		APIResponse(w, r, 200, "User unfollowed", make(map[string]string)) // send data to client side
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
}
