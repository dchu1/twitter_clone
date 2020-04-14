package handlers

import (
	// 	"context"
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"net/http"

	// "net/url"
	// "time"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
)

// PostHandler is the handler for /post. It is used to create posts
func PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.CreatePostRequest{}
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

		//Get user id of the session
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// cookie, err := r.Cookie("sessionId")
		// // if err != nil || cookie.Value != "" {
		// token, _ := url.QueryUnescape(cookie.Value)
		// user, err := AuthClient.GetUserId(ctx, &authpb.AuthToken{Token: token})
		// // }
		// fmt.Println(user)
		user := r.Context().Value("user").(*authpb.UserId)
		_, err = PostServiceClient.CreatePost(r.Context(), &postpb.Post{UserId: user.UserId, Message: reqMessage.Message})
		//err = application.CreatePost(user.UserId, reqMessage.Message)

		if err != nil {
			APIResponse(w, r, http.StatusInternalServerError, "Post not added", make(map[string]string))
			return
		}
		APIResponse(w, r, http.StatusOK, "Post added successfully", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}
