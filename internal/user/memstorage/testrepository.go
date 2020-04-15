package memstorage

import (
	"context"
	"errors"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

type testUserRepository struct {
	userRepo  user.UserRepository
	sleeptime time.Duration
}

const sleeptime = "10s"

// GetUserRepository returns a UserRepository that uses package level storage
func GetTestUserRepository() user.UserRepository {
	ur := GetUserRepository()
	st, _ := time.ParseDuration(sleeptime)
	return &testUserRepository{ur, st}
}

// NewUserRepository reutnrs a UserRepository that uses the given storage
func NewTestUserRepository(ur user.UserRepository) user.UserRepository {
	st, _ := time.ParseDuration(sleeptime)
	return &testUserRepository{ur, st}
}

// CreateUser adds a user to the appropriate data structures
func (testRepo *testUserRepository) CreateUser(ctx context.Context, info *userpb.AccountInformation) (uint64, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.CreateUser(ctx, info)
}

// GetUser creates a copy of the specified user.
func (testRepo *testUserRepository) GetUser(ctx context.Context, userID uint64) (*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetUser(ctx, userID)
}

// GetUsers creates a copy of the specified users.
func (testRepo *testUserRepository) GetUsers(ctx context.Context, userIDs []uint64) ([]*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetUsers(ctx, userIDs)
}

// GetAllUsers returns all users
func (testRepo *testUserRepository) GetAllUsers(ctx context.Context) ([]*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetAllUsers(ctx)
}

// FollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user is following another user
func (testRepo *testUserRepository) FollowUser(ctx context.Context, followingUserID uint64, UserIDToFollow uint64) error {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.FollowUser(ctx, followingUserID, UserIDToFollow)
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (testRepo *testUserRepository) UnFollowUser(ctx context.Context, followingUserID uint64, UserIDToUnfollow uint64) error {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.FollowUser(ctx, followingUserID, UserIDToUnfollow)
}

// GetUserByUsername returns a user object by their username
func (testRepo *testUserRepository) GetUserByUsername(ctx context.Context, email string) (*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetUserByUsername(ctx, email)
}

// GetFollowing returns an array of users that the given user is following
func (testRepo *testUserRepository) GetFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetFollowing(ctx, userId)
}

// GetNotFollowing returns an array of users that the given user is not following
func (testRepo *testUserRepository) GetNotFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.GetNotFollowing(ctx, userId)
}

// DeleteUser removes a user
func (testRepo *testUserRepository) DeleteUser(ctx context.Context, userID uint64) error {
	time.Sleep(testRepo.sleeptime)
	return testRepo.userRepo.DeleteUser(ctx, userID)
}
func (testRepo *testUserRepository) UpdateUserAccountInfo(ctx context.Context, info *userpb.AccountInformation) error {
	time.Sleep(testRepo.sleeptime)
	return errors.New("Feature not implemented")
}
