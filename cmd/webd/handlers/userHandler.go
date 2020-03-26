package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		u := application.GetUsers()
		arr := make([]app.User, len(u))
		for i, v := range u {
			arr[i] = *v
		}
		respMessage := handlermodels.GetUsersResponse{arr}
		b, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	case "POST":
		reqMessage := handlermodels.GetUserRequest{}
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
		u := application.GetUser(reqMessage.UserID)
		b, err = json.Marshal(u)
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

func UserFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		// Get userId from the session cookie
		sess := session.GlobalSessions.SessionQuery(w, r)
		users, err := application.GetFollowing(sess.Get("userId").(uint64))
		respMessage := handlermodels.GetUserFollowingResponse{users}
		body, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

func UserNotFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		// Get userId from the session cookie
		sess := session.GlobalSessions.SessionQuery(w, r)
		users, err := application.GetNotFollowing(sess.Get("userId").(uint64))
		respMessage := handlermodels.GetUserFollowingResponse{users}
		body, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
