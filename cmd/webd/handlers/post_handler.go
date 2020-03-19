package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type PostHandler struct {
	st *test_storage.TestStorage
}

func MakePostHandler(a *test_storage.TestStorage) *PostHandler {
	return &PostHandler{a}
}

func (h *PostHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch req.Method {
	case "GET":
		h.getHandle(head).ServeHTTP(res, req)
	case "PUT":
		h.putHandle().ServeHTTP(res, req)
	default:
		http.Error(res, "Only GET and PUT are allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PostHandler) getHandle(id string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(res, fmt.Sprintf("Invalid post id %q", id), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(res, "Answered by Post Get Handler. Id:%d", id)
	})
}

func (h *PostHandler) putHandle() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Answered by Post Put Handler.")
	})
}
