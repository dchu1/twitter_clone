package handlers

import (
	"fmt"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type FollowHandler struct {
	app *test_storage.TestStorage
}

func MakeFollowHandler(a *test_storage.TestStorage) *FollowHandler {
	return &FollowHandler{a}
}

func (h *FollowHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch req.Method {
	case "PUT":
		h.putHandle(head).ServeHTTP(res, req)
	case "DELETE":
		h.deleteHandle(head).ServeHTTP(res, req)
	default:
		http.Error(res, "Only PUT and DELETE allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FollowHandler) putHandle(id string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Answered by Follower Put Handler.")
	})
}

func (h *FollowHandler) deleteHandle(id string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Answered by Feed Delete Handler.")
	})
}
