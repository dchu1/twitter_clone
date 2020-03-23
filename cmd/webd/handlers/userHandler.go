package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		u := a.GetUsers()
		arr := make([]app.User, len(u))
		for i, v := range u {
			arr[i] = *v
		}
		respMessage := handlermodels.GetUsersResponse{arr}
		b, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	case "POST":
		reqMessage := handlermodels.GetUserRequest{}
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
		u := a.GetUser(reqMessage.UserID)
		b, err = json.Marshal(u)
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
