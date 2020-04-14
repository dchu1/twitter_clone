package handlers

import (
	"encoding/json"
	"net/http"

	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/golang/protobuf/ptypes"
)

// Feed is the Handler for serving request for user's feed
// Gets the user id from the session
func Feed(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// get the user's feed
		user := r.Context().Value("user").(*authpb.UserId)
		// Get the user's obj
		userObj, err := UserServiceClient.GetUser(r.Context(), &userpb.UserId{UserId: user.UserId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Get the user's following list
		followers, err := UserServiceClient.GetFollowing(r.Context(), &userpb.UserId{UserId: user.UserId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Construct an array of userIds to get posts and a map to map userids to user objects
		userMap := make(map[uint64]*userpb.User)
		userMap[user.UserId] = userObj
		tempArr := make([]uint64, 0, len(followers.GetUserList())+1)
		tempArr = append(tempArr, user.UserId)
		for _, v := range followers.GetUserList() {
			tempArr = append(tempArr, v.AccountInformation.UserId)
			userMap[v.AccountInformation.UserId] = v
		}

		// Get posts
		posts, err := PostServiceClient.GetPostsByAuthors(r.Context(), &postpb.UserIDs{UserIDs: tempArr})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//create a reply struct
		respPostArray := make([]handlermodels.Post, len(posts.Posts))
		for i, post := range posts.Posts {
			// construct a post struct
			author := userMap[post.UserId].AccountInformation
			timestamp, _ := ptypes.Timestamp(post.Timestamp)
			authorStruct := handlermodels.Author{UserID: author.UserId, Firstname: author.FirstName, Lastname: author.LastName, Email: author.Email}
			postStruct := handlermodels.Post{Id: post.PostID, Timestamp: timestamp, Message: post.Message, Author: authorStruct}
			respPostArray[i] = postStruct
		}
		respMessage := handlermodels.FeedResponse{respPostArray}
		body, err := json.Marshal(respMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Write(body)
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
