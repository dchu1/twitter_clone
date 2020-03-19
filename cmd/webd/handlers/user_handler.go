package handlers

import (
	"fmt"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type UserHandler struct {
	st *test_storage.TestStorage
}

func MakeUserHandler(a *test_storage.TestStorage) *UserHandler {
	return &UserHandler{a}
}

func (h *UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	// id, err := strconv.Atoi(head)
	// if err != nil {
	//     http.Error(res, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
	//     return
	// }
	if head != "" {
		switch req.Method {
		case "GET":
			h.getHandle(head).ServeHTTP(res, req)
		case "PUT":
			h.putHandle(head).ServeHTTP(res, req)
		default:
			http.Error(res, "Only GET and PUT are allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(res, "Missing Argument", http.StatusBadRequest)
	}
}

func (h *UserHandler) getHandle(name string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Answered by User Get Handler. Name:%s", name)
	})
}

func (h *UserHandler) putHandle(name string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		h.st.AddUser(name)
		fmt.Fprintf(res, "Answered by User Put Handler. Name:%s", name)
	})
}
