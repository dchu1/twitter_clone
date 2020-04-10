package app

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
)

// Service is the interface that provides app methods.
type Service interface {
	CreateUser(context.Context, user.AccountInformation) (uint64, error)
	CreatePost(context.Context, user.AccountInformation, post.Post) (uint64, error)
	GetUser(context.Context, uint64) (*user.User, error)
	GetUsers(context.Context, []uint64) ([]*user.User, error)
	GetFollowing(context.Context, uint64) ([]*user.User, error)
	GetNotFollowing(context.Context, uint64) ([]*user.User, error)
	UpdateUserAccountInfo(context.Context, user.AccountInformation) error
	FollowUser(context.Context, uint64, uint64) error
	UnFollowUser(context.Context, uint64, uint64) error
	DeleteUser(context.Context, uint64) error

	GetFeed(context.Context, uint64) ([]*Post, error)
}

type service struct {
	userRepo user.UserRepository
	postRepo post.PostRepository
}

// Post is a read only struct of posts that includes account information of author
type Post struct {
	PostID    uint64                  `json:"postId,omitempty"`    // This is a unique id. Type might be different depending on how we generate unique ids.
	Timestamp time.Time               `json:"timestamp,omitempty"` // time this post was made
	Message   string                  `json:"message,omitempty"`   // the text of the post
	Author    user.AccountInformation `json:"author,omitempty"`    //id of the user who wrote the post
}

// ByTime implements sort.Interface for []Post based on
// the Timestamp field.
type ByTime []*Post

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return !a[i].Timestamp.Before(a[j].Timestamp) }

func NewService(ur user.UserRepository, pr post.PostRepository) Service {
	return &service{ur, pr}
}

// GetFeed returns a given user's feed. Not sure if this should be in the service layer...
func (s *service) GetFeed(ctx context.Context, userID uint64) ([]*Post, error) {
	retArray := make([]*Post, 0, 100)
	// Get our user
	userObj, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user ids followed
	followed := make([]uint64, len(userObj.Followers))
	// add user to followed
	followed = append(followed, userID)
	for k := range userObj.Followers {
		followed = append(followed, k)
	}
	followedArr, err := s.userRepo.GetUsers(ctx, followed)
	if err != nil {
		return nil, err
	}

	// Get posts
	for _, follower := range followedArr {
		tempArr, err := s.postRepo.GetPosts(ctx, follower.Posts)
		if err != nil {
			return nil, err
		}
		for _, post := range tempArr {
			// construct a new Post
			p := new(Post)
			p.PostID = post.PostID
			p.Timestamp = post.Timestamp
			p.Message = post.Message
			p.Author = follower.AccountInformation
			retArray = append(retArray, p)
		}
	}

	// Sort array
	sort.Sort(ByTime(retArray))
	return retArray, nil
}

func (s *service) CreateUser(ctx context.Context, info user.AccountInformation) (uint64, error) {
	return s.userRepo.CreateUser(ctx, info)
}

func (s *service) CreatePost(ctx context.Context, info user.AccountInformation, p post.Post) (uint64, error) {
	postID, err := s.postRepo.CreatePost(ctx, p)
	if err != nil {
		return 0, err
	}
	err = s.userRepo.AddPost(ctx, info.UserID, postID)
	if err != nil {
		return postID, err
	}
	return postID, nil
}

func (s *service) GetUser(ctx context.Context, userID uint64) (*user.User, error) {
	return s.userRepo.GetUser(ctx, userID)
}
func (s *service) GetUsers(ctx context.Context, userIDs []uint64) ([]*user.User, error) {
	return s.userRepo.GetUsers(ctx, userIDs)
}

func (s *service) GetFollowing(ctx context.Context, userID uint64) ([]*user.User, error) {
	return s.userRepo.GetFollowing(ctx, userID)
}
func (s *service) GetNotFollowing(ctx context.Context, userID uint64) ([]*user.User, error) {
	return s.userRepo.GetNotFollowing(ctx, userID)
}
func (s *service) UpdateUserAccountInfo(ctx context.Context, info user.AccountInformation) error {
	return errors.New("Feature not implemented")
}
func (s *service) FollowUser(ctx context.Context, source uint64, target uint64) error {
	return s.userRepo.FollowUser(ctx, source, target)
}
func (s *service) UnFollowUser(ctx context.Context, source uint64, target uint64) error {
	return s.userRepo.UnFollowUser(ctx, source, target)
}
func (s *service) DeleteUser(ctx context.Context, userID uint64) error {
	return errors.New("Feature not implemented")
}
