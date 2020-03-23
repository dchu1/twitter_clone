package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func FollowCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.FollowRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a.FollowUser(reqMessage.SourceUserId, reqMessage.TargetUserId)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func FollowDestroyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//a.UnFollowUser(nil, nil)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
}
