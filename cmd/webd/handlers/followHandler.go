package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

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
		a.FollowUser(sess.Get("userId").(uint64), reqMessage.TargetUserId)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func FollowDestroyHandler(w http.ResponseWriter, r *http.Request) {
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
		a.UnFollowUser(sess.Get("userId").(uint64), reqMessage.TargetUserId)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
}
