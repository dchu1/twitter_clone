package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

// UserHandler is the handler for /users. It is for getting a list of all users, or a specific user.
func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		u := application.GetUsers()
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
			APIResponse(w, r, http.StatusBadRequest, "Post not added", make(map[string]string))
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			APIResponse(w, r, http.StatusBadRequest, "Post not added", make(map[string]string))
			return
		}
		u := application.GetUser(reqMessage.UserID)
		b, err = json.Marshal(u)
		w.Header().Set("content-type", "application/json")
		w.Write(b)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

// UserFollowingHandler is the handler for /user/following. It is used for getting a list of users
// the current user is following
func UserFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// //Get user id of the session
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// cookie, err := r.Cookie("sessionId")
		// // if err != nil || cookie.Value != "" {
		// token, _ := url.QueryUnescape(cookie.Value)
		// user, err := AuthClient.GetUserId(ctx, &authpb.AuthToken{Token: token})
		// // }
		user := r.Context().Value("user").(*authpb.UserId)
		users, err := UserServiceClient.GetFollowing(r.Context(), &userpb.UserId{UserId: user.UserId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// //users, err := application.GetFollowing(user.UserId)
		// respMessage := handlermodels.GetUserFollowingResponse{users.UserList}
		// body, err := json.Marshal(respMessage)
		body, err := json.Marshal(users.UserList)
		if err != nil {
			APIResponse(w, r, http.StatusInternalServerError, "Cannot get user following list", make(map[string]string))
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}

// UserNotFollowingHandler is the handler for /user/notfollowing. It is used for getting a list of users
// the current user is not following
func UserNotFollowingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING

		// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		// defer cancel()
		// cookie, err := r.Cookie("sessionId")
		// // if err != nil || cookie.Value != "" {
		// token, _ := url.QueryUnescape(cookie.Value)
		// user, err := AuthClient.GetUserId(ctx, &authpb.AuthToken{Token: token})
		// // }

		//users, err := application.GetNotFollowing(user.UserId)
		user := r.Context().Value("user").(*authpb.UserId)
		users, err := UserServiceClient.GetNotFollowing(r.Context(), &userpb.UserId{UserId: user.UserId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// respMessage := handlermodels.GetUserFollowingResponse{users.UserList}
		// body, err := json.Marshal(respMessage)
		body, err := json.Marshal(users.UserList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
