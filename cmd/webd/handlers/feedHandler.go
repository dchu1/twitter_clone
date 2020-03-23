package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sess := session.GlobalSessions.SessionQuery(w, r)
		feed, err := application.GetFeed(sess.Get("userId").(uint64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respPostArray := make([]handlermodels.Post, len(feed))
		for i, post := range feed {
			user := application.GetUser(post.UserID)
			// construct our struct
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
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
