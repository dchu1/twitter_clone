package app

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddUser(t *testing.T) {
	app := MakeApp()
	expected := MakeUser("TestFirst", "TestLast", "test@test.com", "testpass", 0)
	app.AddUser("TestFirst", "TestLast", "test@test.com", "testpass")
	actual := app.GetUsers()[0]
	if expected.Email != actual.Email ||
		expected.FirstName != actual.FirstName ||
		expected.LastName != actual.LastName ||
		expected.id != actual.id {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual.Email, actual.FirstName, actual.LastName, actual.id))
	}
}

func TestFollowUser(t *testing.T) {
	app := MakeApp()
	app.AddUser("FollowerFirst", "FollowerLast", "follower@test.com", "followerpass")
	app.AddUser("FollowingFirst", "FollowingLast", "following@test.com", "followingpass")
	followingUserID := app.users[0].id
	UserIDToFollow := app.users[1].id
	app.FollowUser(followingUserID, UserIDToFollow)
	if app.users[0].following[UserIDToFollow] != UserIDToFollow {
		t.Error("Test Failed Following map not updated properly")
	}
	if app.users[1].followers[followingUserID] != followingUserID {
		t.Error("Test Failed Followers map not updated properly")
	}

}

func TestUnFollowUser(t *testing.T) {
	app := MakeApp()
	app.AddUser("User0First", "User0Last", "User0@test.com", "User0pass")
	app.AddUser("User1First", "User1Last", "User1@test.com", "User1gpass")
	app.AddUser("User2First", "User2Last", "User2@test.com", "User2gpass")

	User0FollowingList := map[uint64]uint64{2: 2}
	User1FollowerList := make(map[uint64]uint64)
	User2FollowerList := map[uint64]uint64{0: 0}

	User0ID := app.users[0].id
	User1ID := app.users[1].id
	User2ID := app.users[2].id
	app.FollowUser(User0ID, User1ID)
	app.FollowUser(User0ID, User2ID)
	app.UnFollowUser(User0ID, User1ID)

	if reflect.DeepEqual(app.users[0].following, User0FollowingList) == false {
		t.Error("Test Failed Following map not updated properly")
	}
	if reflect.DeepEqual(app.users[1].followers, User1FollowerList) == false {
		t.Error("Test Failed Followers map not updated properly for User 1")
	}
	if reflect.DeepEqual(app.users[2].followers, User2FollowerList) == false {
		t.Error("Test Failed Followers map not updated properly User 2")
	}
}

func TestCreatePost(t *testing.T) {
	app := MakeApp()
	app.AddUser("TestFirst", "TestLast", "Test@test.com", "Testpass")
	userID := app.users[0].id
	app.CreatePost(userID, "Test Message")
	userPost := app.users[0].post[0]
	appPost := app.posts[0]

	if reflect.DeepEqual(userPost, appPost) == false {
		t.Error("Test Failed User struct not updated properly for the added post")
	}
	if userPost.Message != "Test Message" {
		t.Error("Test Failed User struct not updated properly for Post message ")
	}
	posts, _ := app.GetFeed(userID)
	if reflect.DeepEqual(appPost, posts[0]) == false {
		t.Error("Test Failed Posts struct not in sync")
	}
}

func TestGetFeed(t *testing.T) {
	app := MakeApp()
	app.AddUser("TestFirst", "TestLast", "Test@test.com", "Testpass")
	userID := app.users[0].id
	app.CreatePost(userID, "Test Message")
	userPost := app.users[0].post[0]

	posts, _ := app.GetFeed(userID)

	if reflect.DeepEqual(userPost, posts[0]) == false {
		t.Error("Test Failed added post not returned in the list of feeds")
	}
}
