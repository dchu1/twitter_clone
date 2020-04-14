package user

import (
	"context"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

// AccountInformation represents account information
// type AccountInformation struct {
// 	FirstName string `json:"firstname,omitempty"`
// 	LastName  string `json:"lastname,omitempty"`
// 	Email     string `json:"email,omitempty"`
// 	UserID    uint64 `json:"userId"`
// }

// type User struct {
// 	AccountInformation AccountInformation
// 	Following          map[uint64]struct{}
// 	Followers          map[uint64]struct{}
// 	//Posts              []uint64
// }

type UserRepository interface {
	CreateUser(context.Context, *pb.AccountInformation) (uint64, error)
	GetUser(context.Context, uint64) (*pb.User, error)
	GetUsers(context.Context, []uint64) ([]*pb.User, error)
	GetAllUsers(context.Context) ([]*pb.User, error)
	GetUserByUsername(context.Context, string) (*pb.User, error)
	GetFollowing(context.Context, uint64) ([]*pb.User, error)
	GetNotFollowing(context.Context, uint64) ([]*pb.User, error)
	UpdateUserAccountInfo(context.Context, *pb.AccountInformation) error
	FollowUser(context.Context, uint64, uint64) error
	UnFollowUser(context.Context, uint64, uint64) error
	DeleteUser(context.Context, uint64) error
}

// copyFollowMap makes a deep copy of a user's following or followed map
func copyFollowMap(m map[uint64]struct{}) map[uint64]struct{} {
	cp := make(map[uint64]struct{})
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// func (user *User) Clone() *User {
// 	retUser := user
// 	retUser.Following = copyFollowMap(retUser.Following)
// 	retUser.Followers = copyFollowMap(retUser.Followers)
// 	return retUser
// }

// func NewUser(info AccountInformation) *User {
// 	return &User{AccountInformation: info, Following: make(map[uint64]struct{}), Followers: make(map[uint64]struct{})}
// }
