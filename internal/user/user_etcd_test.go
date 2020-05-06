package user_test

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/storage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/storage/etcd"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"go.etcd.io/etcd/clientv3"
)

func CleanUpEtcd(cli *clientv3.Client) {
	cli.Delete(context.Background(), "/User/", clientv3.WithPrefix())
	cli.Put(context.Background(), "/NextUserId", "1")
}

func TestAddUserEtcd(t *testing.T) {
	//userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	expected_user := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}

	userId, err := userApp.CreateUser(context.Background(), &expected_user)
	if err != nil {
		t.Error(err.Error())
		return
	}
	actual_user, _ := userApp.GetUser(context.Background(), userId)
	if expected_user.Email != actual_user.AccountInformation.Email ||
		expected_user.FirstName != actual_user.AccountInformation.FirstName ||
		expected_user.LastName != actual_user.AccountInformation.LastName {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual_user.AccountInformation.Email, actual_user.AccountInformation.FirstName, actual_user.AccountInformation.LastName, actual_user.AccountInformation.UserId))
	}
}

func TestGetUsersEtcd(t *testing.T) {
	const numUsers = 10
	//userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)
	expectedUsers := make(map[uint64]*userpb.AccountInformation)
	retUsers := make(map[uint64]*userpb.AccountInformation)
	userIds := make([]uint64, numUsers)
	for i := 1; i <= numUsers; i++ {
		user := userpb.AccountInformation{UserId: uint64(i), FirstName: "test" + strconv.Itoa(i), LastName: "test" + strconv.Itoa(i), Email: strconv.Itoa(i) + "@test.edu"}
		userIdpb, err := userApp.CreateUser(context.Background(), &user)
		if err != nil {
			t.Error(err.Error())
		}
		expectedUsers[userIdpb.UserId] = &user
		userIds[i-1] = userIdpb.UserId
	}
	ret, _ := userApp.GetUsers(context.Background(), &userpb.UserIds{UserIds: userIds})
	if len(ret.UserList) != numUsers {
		t.Error(fmt.Sprintf("Unexpected number of users. Expected:%d, Got:%d, %v\n", numUsers, len(ret.UserList), ret.UserList))
		return
	}

	for _, user := range ret.UserList {
		retUsers[user.AccountInformation.UserId] = user.AccountInformation
	}
	for _, actual_user := range retUsers {
		expected_user := expectedUsers[actual_user.UserId]
		if expected_user.Email != actual_user.Email ||
			expected_user.FirstName != actual_user.FirstName ||
			expected_user.LastName != actual_user.LastName ||
			expected_user.UserId != actual_user.UserId {
			t.Error(fmt.Sprintf("Test Failed: Got: %+v, Expected: %+v\n", actual_user, expected_user))
		}
	}
}

func TestFollowUserEtcd(t *testing.T) {
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
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

func TestUnFollowUserEtcd(t *testing.T) {
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
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

func TestConcurrentAddUserEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	for user := 0; user < numUsers; user++ {
		go func(user int) {
			defer wg.Done()
			firstName := "TestFirstName" + strconv.Itoa(user)
			lastName := "TestLastName" + strconv.Itoa(user)
			email := "TestEmail" + strconv.Itoa(user)
			_, err := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: firstName, LastName: lastName, Email: email})
			if err != nil {
				t.Error(err.Error())
			}
		}(user)
	}
	wg.Wait()
	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList) != numUsers {
		t.Error("All users not added in the struct")
	}
}

func TestConcurrentFollowEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	//Create Users
	for user := 1; user <= numUsers; user++ {
		go func(user int) {
			defer wg.Done()
			firstName := "TestFirstName" + strconv.Itoa(user)
			lastName := "TestLastName" + strconv.Itoa(user)
			email := "TestEmail" + strconv.Itoa(user)
			_, err := userApp.CreateUser(context.Background(), &userpb.AccountInformation{FirstName: firstName, LastName: lastName, Email: email})
			if err != nil {
				t.Error(err.Error())
				return
			}
		}(user)
	}
	wg.Wait()
	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList) != numUsers {
		t.Error("All users not created")
		return
	}
	wg.Add(numUsers)

	// Have them all follow each other
	for i := 1; i <= numUsers; i++ {
		go func(userId uint64) {
			defer wg.Done()
			for k := 1; k <= numUsers; k++ {
				if userId != uint64(k) {
					_, err := userApp.FollowUser(context.Background(), &userpb.FollowRequest{UserId: userId, FollowUserId: uint64(k)})
					if err != nil {
						t.Error(err.Error())
					}
				}
			}
		}(uint64(i))
	}

	wg.Wait()

	for i := 1; i <= numUsers; i++ {
		user, _ := userApp.GetUser(context.Background(), &userpb.UserId{UserId: uint64(i)})
		if len(user.Following) != (numUsers - 1) {
			t.Error(fmt.Sprintf("Following map of user %d only has length %d", i, len(user.Following)))
		}
	}
}

func TestContextTimeoutAddUserEtcd(t *testing.T) {
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	// Mock repository with 10 seconds delay for accessing database
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	testUserRepo := storage.NewTestUserRepository(userRepo)
	userApp := user.GetUserServiceServer(&testUserRepo)

	expected_user := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	userApp.CreateUser(ctx, &expected_user)

	users, _ := userApp.GetAllUsers(context.Background(), nil)
	if len(users.UserList) > 0 {
		t.Error(fmt.Sprintf("Test Failed: User added even when context was cancelled"))
	}
	cancel()
}

func TestContextTimeoutFollowUserEtcd(t *testing.T) {

	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	// Mock repository with 10 seconds delay for accessing database
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	testUserRepo := storage.NewTestUserRepository(userRepo)
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

func TestContextTimeoutUnFollowUserEtcd(t *testing.T) {

	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	// Mock repository with 10 seconds delay for accessing database
	userStorage, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	CleanUpEtcd(userStorage)

	userRepo := etcd.NewUserRepository(userStorage)
	testUserRepo := storage.NewTestUserRepository(userRepo)
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
