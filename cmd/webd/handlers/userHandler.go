package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

// UserHandler is the handler for /users. It is for getting a list of all users, or a specific user.
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

// UserFollowingHandler is the handler for /user/following. It is used for getting a list of users
// the current user is following
func UserFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get the session from the context
		sess, ok := session.FromContext(r.Context())
		if !ok {
			http.Error(w, "Context has no session", http.StatusInternalServerError)
			return
		}
		users, err := application.GetFollowing(sess.Get("userId").(uint64))
		respMessage := handlermodels.GetUserFollowingResponse{users}
		body, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

// UserNotFollowingHandler is the handler for /user/notfollowing. It is used for getting a list of users
// the current user is not following
func UserNotFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		// Get the session from the context
		sess, ok := session.FromContext(r.Context())
		if !ok {
			http.Error(w, "Context has no session", http.StatusInternalServerError)
			return
		}
		users, err := application.GetNotFollowing(sess.Get("userId").(uint64))
		respMessage := handlermodels.GetUserFollowingResponse{users}
		body, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
