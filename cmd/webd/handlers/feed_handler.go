package handlers

import (
	"fmt"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type FeedHandler struct {
	st *test_storage.TestStorage
}

func MakeFeedHandler(a *test_storage.TestStorage) *FeedHandler {
	return &FeedHandler{a}
}

func (h *FeedHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		h.getHandle().ServeHTTP(res, req)
	default:
		http.Error(res, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FeedHandler) getHandle() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Answered by Feed Get Handler.")
	})
}
