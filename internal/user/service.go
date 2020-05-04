package user

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

type userServiceServer struct {
	userRepo UserRepository
	pb.UnimplementedUserServiceServer
}

func (s *userServiceServer) CreateUser(ctx context.Context, info *pb.AccountInformation) (*pb.UserId, error) {
	//Check whether user already exists
	userObj, _ := s.userRepo.GetUserByUsername(ctx, info.Email)
	if userObj != nil {
		return &pb.UserId{}, errors.New("duplicate email")
	}

	// Create user object
	user := new(pb.User)
	user.AccountInformation = info
	id, err := s.userRepo.NextUserId()
	if err != nil {
		return &pb.UserId{}, fmt.Errorf("could not get user id\n%s", err)
	}
	user.AccountInformation.UserId = id
	user.Followers = make(map[uint64]uint64)
	user.Following = make(map[uint64]uint64)

	uid, err := s.userRepo.CreateUser(ctx, user)
	return &pb.UserId{UserId: uid}, err
}
func (s *userServiceServer) GetUser(ctx context.Context, req *pb.UserId) (*pb.User, error) {
	return s.userRepo.GetUser(ctx, req.GetUserId())
}
func (s *userServiceServer) GetUsers(ctx context.Context, req *pb.UserIds) (*pb.UserList, error) {
	if len(req.GetUserIds()) == 0 {
		return new(pb.UserList), nil
	}
	users, err := s.userRepo.GetUsers(ctx, req.GetUserIds())
	return &pb.UserList{UserList: users}, err
}
func (s *userServiceServer) GetAllUsers(ctx context.Context, req *pb.Void) (*pb.UserList, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	return &pb.UserList{UserList: users}, err
}
func (s *userServiceServer) GetFollowing(ctx context.Context, req *pb.UserId) (*pb.UserList, error) {
	users, err := s.userRepo.GetFollowing(ctx, req.GetUserId())
	return &pb.UserList{UserList: users}, err
}
func (s *userServiceServer) GetNotFollowing(ctx context.Context, req *pb.UserId) (*pb.UserList, error) {
	users, err := s.userRepo.GetNotFollowing(ctx, req.GetUserId())
	return &pb.UserList{UserList: users}, err
}
func (s *userServiceServer) FollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.Void, error) {
	if req.GetUserId() == req.GetFollowUserId() {
		return &pb.Void{}, errors.New("duplicate user ids")
	}
	return &pb.Void{}, s.userRepo.FollowUser(ctx, req.GetUserId(), req.GetFollowUserId())
}
func (s *userServiceServer) UnFollowUser(ctx context.Context, req *pb.UnFollowRequest) (*pb.Void, error) {
	if req.GetUserId() == req.GetFollowUserId() {
		return &pb.Void{}, errors.New("duplicate user ids")
	}
	return &pb.Void{}, s.userRepo.UnFollowUser(ctx, req.GetUserId(), req.GetFollowUserId())
}
func (s *userServiceServer) GetUserIdByUsername(ctx context.Context, req *pb.UserName) (*pb.UserId, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.GetEmail())
	return &pb.UserId{UserId: user.AccountInformation.UserId}, err
}

// GetUserServiceServer returns a grpc Server for the user service using the provided UserRepository
func GetUserServiceServer(ur *UserRepository) *userServiceServer {
	return &userServiceServer{userRepo: *ur}
}
