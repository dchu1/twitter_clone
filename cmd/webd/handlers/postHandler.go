package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

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

		// Get the session from the context
		sess, ok := session.FromContext(r.Context())
		if !ok {
			http.Error(w, "Context has no session", http.StatusInternalServerError)
			return
		}
		err = application.CreatePost(sess.Get("userId").(uint64), reqMessage.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		APIResponse(w, r, http.StatusOK, "Post added successfully", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}
