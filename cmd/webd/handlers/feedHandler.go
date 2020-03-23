package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sess := session.GlobalSessions.SessionQuery(w, r)
		feed, err := a.GetFeed(sess.Get("userId").(uint64))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respPostArray := make([]handlermodels.Post, len(feed))
		for i, v := range feed {
			u := a.GetUser(v.UserID)
			// construct our struct
			a := handlermodels.Author{UserID: u.Id, Firstname: u.FirstName, Lastname: u.LastName, Email: u.Email}
			p := handlermodels.Post{Id: v.Id, Timestamp: v.Timestamp, Message: v.Message, Author: a}
			respPostArray[i] = p
		}
		respMessage := handlermodels.FeedResponse{respPostArray}
		b, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	case "POST":
		reqMessage := handlermodels.FeedRequest{}
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
		feed, err := a.GetFeed(reqMessage.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respPostArray := make([]handlermodels.Post, len(feed))
		for i, v := range feed {
			u := a.GetUser(v.UserID)
			// construct our struct
			a := handlermodels.Author{UserID: u.Id, Firstname: u.FirstName, Lastname: u.LastName, Email: u.Email}
			p := handlermodels.Post{Id: v.Id, Timestamp: v.Timestamp, Message: v.Message, Author: a}
			respPostArray[i] = p
		}
		respMessage := handlermodels.FeedResponse{respPostArray}
		b, err = json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(b)
		//handler.APIResponse(w, r, http.StatusOK, "Added successful;y", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}
