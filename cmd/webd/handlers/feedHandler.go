package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

// Feed is the Handler for serving request for user's feed
// Gets the user id from the session
func Feed(w http.ResponseWriter, r *http.Request) {


	switch r.Method {
	case "GET":
		// get the session from the cookie
		sess := session.GlobalSessions.SessionQuery(w, r)

		// get the user's feed
		feed, err := application.GetFeed(sess.Get("userId").(uint64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// create a reply struct
		respPostArray := make([]handlermodels.Post, len(feed))
		for i, post := range feed {
			user := application.GetUser(post.UserID)
			// construct a post struct
			authorStruct := handlermodels.Author{UserID: user.Id, Firstname: user.FirstName, Lastname: user.LastName, Email: user.Email}
			postStruct := handlermodels.Post{Id: post.Id, Timestamp: post.Timestamp, Message: post.Message, Author: authorStruct}
			respPostArray[i] = postStruct
		}
		respMessage := handlermodels.FeedResponse{respPostArray}
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
