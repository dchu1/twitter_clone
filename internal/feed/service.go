package feed

import (
	"context"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/feed/feedpb"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

type feedServiceServer struct {
	PostServiceClient postpb.PostServiceClient
	UserServiceClient userpb.UserServiceClient
	feedpb.UnimplementedFeedServiceServer
}

func (s *feedServiceServer) GetFeed(ctx context.Context, req *feedpb.UserId) (*feedpb.FeedMessage, error) {
	// TODO Convert the feedb.AccountInformation to userpb.AccountInformation and same for post
	// Ideally I would use the same object, but I don't know how to import protoc files from
	// other protoc files...

	errorchan := make(chan error, 1)
	resultchan := make(chan *feedpb.FeedMessage, 1)
	go func() {
		userObj, err := s.UserServiceClient.GetUser(ctx, &userpb.UserId{UserId: req.UserId})
		if err != nil {
			errorchan <- err
			return
		}
		// Get the user's following list
		followers, err := s.UserServiceClient.GetFollowing(ctx, &userpb.UserId{UserId: req.UserId})
		if err != nil {
			errorchan <- err
			return
		}
		// Construct an array of userIds to get posts and a map to map userids to user objects
		userMap := make(map[uint64]*userpb.User)
		userMap[req.UserId] = userObj
		tempArr := make([]uint64, 0, len(followers.GetUserList())+1)
		tempArr = append(tempArr, req.UserId)
		for _, v := range followers.GetUserList() {
			tempArr = append(tempArr, v.AccountInformation.UserId)
			userMap[v.AccountInformation.UserId] = v
		}

		// Get posts
		posts, err := s.PostServiceClient.GetPostsByAuthors(ctx, &postpb.UserIDs{UserIDs: tempArr})
		if err != nil {
			errorchan <- err
			return
		}

		//create a reply struct
		respPostArray := make([]*feedpb.Post, len(posts.Posts))
		for i, post := range posts.Posts {
			// construct a post struct
			userObj = userMap[post.UserId]
			authorMessageObj := &feedpb.AccountInformation{FirstName: userObj.AccountInformation.FirstName, LastName: userObj.AccountInformation.LastName, Email: userObj.AccountInformation.Email, UserId: userObj.AccountInformation.UserId}
			postMessageObj := feedpb.Post{PostID: post.PostID, Timestamp: post.Timestamp, Message: post.Message, Author: authorMessageObj}
			respPostArray[i] = &postMessageObj
		}

		resultchan <- &feedpb.FeedMessage{Posts: respPostArray}
	}()

	select {
	case res := <-resultchan:
		return res, nil
	case err := <-errorchan:
		return new(feedpb.FeedMessage), err
	case <-ctx.Done():
		return new(feedpb.FeedMessage), ctx.Err()
	}

}

// NewFeedServiceServer returns a grpc Server for the feed service using the provided clients
func NewFeedServiceServer(psc postpb.PostServiceClient, usc userpb.UserServiceClient) *feedServiceServer {
	return &feedServiceServer{PostServiceClient: psc, UserServiceClient: usc}
}
