package handlers

import (
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type LoginHandler struct {
	st *test_storage.TestStorage
}

func (h *LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
