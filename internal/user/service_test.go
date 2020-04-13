package user_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/memstorage"
)

func TestAddUser(t *testing.T) {
	storage := memstorage.NewUserStorage()
	app := user.NewService(memstorage.NewUserRepository(storage))
	expected := user.AccountInformation{"TestFirst", "TestLast", "test@test.com", 0}
	userID, _ := app.CreateUser(nil, expected)
	actual, _ := app.GetUser(nil, userID)
	if expected.Email != actual.AccountInformation.Email ||
		expected.FirstName != actual.AccountInformation.FirstName ||
		expected.LastName != actual.AccountInformation.LastName ||
		expected.UserID != actual.AccountInformation.UserID {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual.AccountInformation.Email, actual.AccountInformation.FirstName, actual.AccountInformation.LastName, actual.AccountInformation.UserID))
	}
}

func TestFollowUser(t *testing.T) {
	storage := memstorage.NewUserStorage()
	app := user.NewService(memstorage.NewUserRepository(storage))
	followingUserID, _ := app.CreateUser(nil, user.AccountInformation{"FollowerFirst", "FollowerLast", "follower@test.com", 0})
	UserIDToFollow, _ := app.CreateUser(nil, user.AccountInformation{"FollowingFirst", "FollowingLast", "following@test.com", 0})
	err := app.FollowUser(nil, followingUserID, UserIDToFollow)
	if err != nil {
		t.Error(err.Error())
	}
	follower, _ := app.GetUser(nil, followingUserID)
	followed, _ := app.GetUser(nil, UserIDToFollow)
	if _, exists := follower.Following[UserIDToFollow]; !exists {
		t.Error(fmt.Sprintf("Test Failed Following map not updated properly: %v", follower.Following))
	}
	if _, exists := followed.Followers[followingUserID]; !exists {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly: %v", followed.Followers))
	}
}

func TestUnFollowUser(t *testing.T) {
	storage := memstorage.NewUserStorage()
	app := user.NewService(memstorage.NewUserRepository(storage))
	User0ID, _ := app.CreateUser(nil, user.AccountInformation{"User0First", "User0Last", "User0@test.com", 0})
	User1ID, _ := app.CreateUser(nil, user.AccountInformation{"User1First", "User1Last", "User1@test.com", 0})
	User2ID, _ := app.CreateUser(nil, user.AccountInformation{"User2First", "User2Last", "User2@test.com", 0})

	User0FollowingList := map[uint64]struct{}{2: struct{}{}}
	User1FollowerList := make(map[uint64]struct{})
	User2FollowerList := map[uint64]struct{}{0: struct{}{}}

	app.FollowUser(nil, User0ID, User1ID)
	app.FollowUser(nil, User0ID, User2ID)
	app.UnFollowUser(nil, User0ID, User1ID)

	u1, _ := app.GetUser(nil, User0ID)
	u2, _ := app.GetUser(nil, User1ID)
	u3, _ := app.GetUser(nil, User2ID)
	if reflect.DeepEqual(u1.Following, User0FollowingList) == false {
		t.Error("Test Failed Following map not updated properly")
	}
	if reflect.DeepEqual(u2.Followers, User1FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 1: %v", u2.Followers))
	}
	if reflect.DeepEqual(u3.Followers, User2FollowerList) == false {
		t.Error("Test Failed Followers map not updated properly User 2")
	}
}
