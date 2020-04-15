package user_test

import (
	"context"
	"fmt"
	"sync"
	"reflect"
	"testing"
	"strconv"
	"time"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/memstorage"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

func TestAddUser(t *testing.T) {
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	expected_user := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}

	userId, _ := userApp.CreateUser(context.Background(), &expected_user)
	actual_user, _ := userApp.GetUser(context.Background(), userId)
	if expected_user.Email != actual_user.AccountInformation.Email ||
		expected_user.FirstName != actual_user.AccountInformation.FirstName ||
		expected_user.LastName != actual_user.AccountInformation.LastName ||
		expected_user.UserId != actual_user.AccountInformation.UserId {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual_user.AccountInformation.Email, actual_user.AccountInformation.FirstName, actual_user.AccountInformation.LastName, actual_user.AccountInformation.UserId))
	}
}

func TestFollowUser(t *testing.T) {
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	followingUserID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "FollowerFirst", LastName: "FollowerLast", Email: "follower@test.com"})
	UserIDToFollow, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "FollowingFirst", LastName: "FollowingLast", Email: "following@test.com"})
	_, err := userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: followingUserID.UserId, FollowUserId: UserIDToFollow.UserId})
	if err != nil {
		t.Error(err.Error())
	}
	follower, _ := userApp.GetUser(context.Background(), followingUserID)
	followed, _ := userApp.GetUser(context.Background(), UserIDToFollow)
	if _, exists := follower.Following[UserIDToFollow.UserId]; !exists {
		t.Error(fmt.Sprintf("Test Failed Following map not updated properly: %v", follower.Following))
	}
	if _, exists := followed.Followers[followingUserID.UserId]; !exists {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly: %v", followed.Followers))
	}
}

func TestUnFollowUser(t *testing.T) {
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)	
	userApp := user.GetUserServiceServer(&userRepo)

	User0ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User0First", LastName: "User0Last", Email: "User0@test.com"})
	User1ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User1First", LastName: "User1Last", Email: "User1@test.com"})
	User2ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User2First", LastName: "User2Last", Email: "User2@test.com"})

	// User0FollowingList := map[uint64]struct{}{2: struct{}{}}
	User0FollowingList := map[uint64]uint64{3: 3}
	User1FollowerList := make(map[uint64]uint64)
	User2FollowerList := make(map[uint64]uint64)

	userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User1ID.UserId})
	userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User2ID.UserId})
	userApp.UnFollowUser(context.Background(), &userpb.UnFollowRequest{UserId: User0ID.UserId, FollowUserId: User1ID.UserId})

	u0, _ := userApp.GetUser(context.Background(), User0ID)
	u1, _ := userApp.GetUser(context.Background(), User1ID)
	u2, _ := userApp.GetUser(context.Background(), User2ID)
	if reflect.DeepEqual(u0.Following, User0FollowingList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 0: %v", u0.Following))
	}
	if reflect.DeepEqual(u1.Following, User1FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 1: %v", u1.Following))
	}
	if reflect.DeepEqual(u2.Following, User2FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 2: %v", u2.Following))
	}
}

func TestContextAddUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	testUserRepo := memstorage.NewTestUserRepository(userRepo)
	userApp := user.GetUserServiceServer(&testUserRepo)
	
	expected_user := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	cancel()
	userApp.CreateUser(ctx, &expected_user)

	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList)>0{
		t.Error(fmt.Sprintf("Test Failed: User added even when context was cancelled"))
	}
}

func TestConcurrentAddUser(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	for user := 0; user < numUsers; user++ {
		go func(user int) {
			defer wg.Done()
			firstName := "TestFirstName" + strconv.Itoa(user)
			lastName := "TestLastName" + strconv.Itoa(user)
			email := "TestEmail" + strconv.Itoa(user)
			userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: firstName, LastName: lastName, Email: email})
		}(user)
	}
	wg.Wait()
	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList) != 100 {
		t.Error("All users not added in the struct")
	}
}

func TestConcurrentFollow(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	//Create Users
	for user := 0; user < numUsers; user++ {
		go func(user int) {
			defer wg.Done()
			firstName := "TestFirstName" + strconv.Itoa(user)
			lastName := "TestLastName" + strconv.Itoa(user)
			email := "TestEmail" + strconv.Itoa(user)
			userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: firstName, LastName: lastName, Email: email})
		}(user)
	}
	wg.Wait()
	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList) != 100 {
		t.Error("All users not added in the struct")
	}
	wg.Add(numUsers)

	// Have them all follow each other and then make a post
	for i := 0; i < numUsers; i++ {
		go func(userId uint64) {
			defer wg.Done()
			for k := 0; k < numUsers; k++ {
				userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: userId, FollowUserId:uint64(k)})
			}
		}(uint64(i))
	}

	wg.Wait()

	for i := 1; i < numUsers; i++ {
		user, _ := userApp.GetUser(context.Background(), &userpb.UserId{UserId:uint64(i)})
		if len(user.Following)!= (numUsers-1){
			t.Error(fmt.Sprintf("Following map of user %d ",i))
		}
	}
}

func TestContextTimeoutAddUser(t *testing.T) {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	
	// Mock repository with 10 seconds delay for accessing database
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	testUserRepo := memstorage.NewTestUserRepository(userRepo)
	userApp := user.GetUserServiceServer(&testUserRepo)
	
	expected_user := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	userApp.CreateUser(ctx, &expected_user)

	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList)>0{
		t.Error(fmt.Sprintf("Test Failed: User added even when context was cancelled"))
	}
	cancel()
}


func TestContextTimeoutFollowUser(t *testing.T) {

	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	// Mock repository with 10 seconds delay for accessing database
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	testUserRepo := memstorage.NewTestUserRepository(userRepo)
	userApp := user.GetUserServiceServer(&testUserRepo)

	User0ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User0First", LastName: "User0Last", Email: "User0@test.com"})
	User1ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User1First", LastName: "User1Last", Email: "User1@test.com"})
	User2ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User2First", LastName: "User2Last", Email: "User2@test.com"})

	User0FollowingList := map[uint64]uint64{3: 3}
	User1FollowerList := make(map[uint64]uint64)
	User2FollowerList := make(map[uint64]uint64)

	//User follow should be unsuccessful because of timeout
	userApp.FollowUser(ctx, &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User1ID.UserId})
	userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User2ID.UserId})	

	u0, _ := userApp.GetUser(context.Background(), User0ID)
	u1, _ := userApp.GetUser(context.Background(), User1ID)
	u2, _ := userApp.GetUser(context.Background(), User2ID)
	if reflect.DeepEqual(u0.Following, User0FollowingList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 0: %v", u0.Following))
	}
	if reflect.DeepEqual(u1.Following, User1FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 1: %v", u1.Following))
	}
	if reflect.DeepEqual(u2.Following, User2FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 2: %v", u2.Following))
	}
}

func TestContextTimeoutUnFollowUser(t *testing.T) {

	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	// Mock repository with 10 seconds delay for accessing database
	userStorage := memstorage.NewUserStorage()
	userRepo := memstorage.NewUserRepository(userStorage)
	testUserRepo := memstorage.NewTestUserRepository(userRepo)
	userApp := user.GetUserServiceServer(&testUserRepo)

	User0ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User0First", LastName: "User0Last", Email: "User0@test.com"})
	User1ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User1First", LastName: "User1Last", Email: "User1@test.com"})
	User2ID, _ := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: "User2First", LastName: "User2Last", Email: "User2@test.com"})

	User0FollowingList := map[uint64]uint64{3: 3}
	User1FollowerList := make(map[uint64]uint64)
	User2FollowerList := make(map[uint64]uint64)

	
	userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User1ID.UserId})
	userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: User0ID.UserId, FollowUserId: User2ID.UserId})	
	//User unfollow should be unsuccessful because of timeout
	userApp.UnFollowUser(ctx, &userpb.UnFollowRequest{UserId: User0ID.UserId, FollowUserId: User1ID.UserId})

	u0, _ := userApp.GetUser(context.Background(), User0ID)
	u1, _ := userApp.GetUser(context.Background(), User1ID)
	u2, _ := userApp.GetUser(context.Background(), User2ID)
	if reflect.DeepEqual(u0.Following, User0FollowingList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 0: %v", u0.Following))
	}
	if reflect.DeepEqual(u1.Following, User1FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 1: %v", u1.Following))
	}
	if reflect.DeepEqual(u2.Following, User2FollowerList) == false {
		t.Error(fmt.Sprintf("Test Failed Followers map not updated properly for User 2: %v", u2.Following))
	}
}