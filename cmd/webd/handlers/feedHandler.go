package handlers

import (
	// "context"
	"encoding/json"
	"net/http"

	// "net/url"
	// "time"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

// Feed is the Handler for serving request for user's feed
// Gets the user id from the session
func Feed(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//Get user id of the session
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// cookie, err := r.Cookie("sessionId")
		// // if err != nil || cookie.Value != "" {
		// token, _ := url.QueryUnescape(cookie.Value)
		// user, err := AuthClient.GetUserId(ctx, &authpb.AuthToken{Token: token})
		// // }

		// get the user's feed
		user := r.Context().Value("user").(*authpb.UserId)
		feed, err := application.GetFeed(user.UserId)
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
