package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

		user := r.Context().Value("user").(*authpb.UserId)
		_, err = PostServiceClient.CreatePost(r.Context(), &postpb.Post{UserId: user.UserId, Message: reqMessage.Message})

		if err != nil {
			APIResponse(w, r, http.StatusInternalServerError, "Post not added", make(map[string]string))
			return
		}
		APIResponse(w, r, http.StatusOK, "Post added successfully", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}
