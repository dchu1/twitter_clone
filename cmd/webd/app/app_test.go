package app

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

func TestAddUser(t *testing.T) {
	app := MakeApp()
	expected := MakeUser("TestFirst", "TestLast", "test@test.com", 0)
	app.AddUser("TestFirst", "TestLast", "test@test.com", "testpass")
	actual := app.GetUsers()[0]
	if expected.Email != actual.Email ||
		expected.FirstName != actual.FirstName ||
		expected.LastName != actual.LastName ||
		expected.Id != actual.Id {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual.Email, actual.FirstName, actual.LastName, actual.Id))
	}
}

func TestFollowUser(t *testing.T) {
	app := MakeApp()
	app.AddUser("FollowerFirst", "FollowerLast", "follower@test.com", "followerpass")
	app.AddUser("FollowingFirst", "FollowingLast", "following@test.com", "followingpass")
	followingUserID := app.users[0].Id
	UserIDToFollow := app.users[1].Id
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

	User0ID := app.users[0].Id
	User1ID := app.users[1].Id
	User2ID := app.users[2].Id
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
	userID := app.users[0].Id
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
	if reflect.DeepEqual(*appPost, posts[0]) == false {
		t.Error("Test Failed Posts struct not in sync")
	}
}

func TestGetFeed(t *testing.T) {
	app := MakeApp()
	app.AddUser("TestFirst", "TestLast", "Test@test.com", "Testpass")
	userID := app.users[0].Id
	app.CreatePost(userID, "Test Message")
	userPost := app.users[0].post[0]

	posts, _ := app.GetFeed(userID)

	if reflect.DeepEqual(*userPost, posts[0]) == false {
		t.Error("Test Failed added post not returned in the list of feeds")
	}
}

func TestConcurrentFollow(t *testing.T) {
	var wg sync.WaitGroup
	rand.Seed(42)
	numUsers := 10
	wg.Add(numUsers)
	app := MakeApp()
	// Create users
	for i := 0; i < numUsers; i++ {
		app.AddUser(strconv.Itoa(i), strconv.Itoa(i), strconv.Itoa(i), strconv.Itoa(i))
	}

	// Have them all follow each other and then make a post
	for i := 0; i < numUsers; i++ {
		go func(userId uint64) {
			defer wg.Done()
			for k := 0; k < numUsers; k++ {
				app.FollowUser(userId, uint64(k))
			}
			// sleep for some random amount of time between 0 and 5 seconds
			//time.Sleep(time.Second * time.Duration(rand.Intn(20)))
			app.CreatePost(userId, strconv.FormatUint(userId, 10))
		}(uint64(i))
	}

	wg.Wait()

	if len(app.posts) != numUsers {
		t.Error(fmt.Sprintf("Incorrect # of posts. Expected %d, found %d", numUsers, len(app.posts)))
	}

	// Get each users feed
	feeds := make([][]Post, numUsers)
	for i := 0; i < numUsers; i++ {
		var err error
		feeds[i], err = app.GetFeed(uint64(i))
		if err != nil {
			t.Error(err.Error())
		}
		if len(feeds[i]) != numUsers {
			t.Error(fmt.Sprintf("Not enough posts in user %d feed. Expected %d, found %d", i, numUsers, len(feeds[i])))
		}
	}

	// Print out user 0 feed
	fmt.Printf("User 0 Feed: ")
	for i := 0; i < numUsers; i++ {
		fmt.Printf("%s, ", feeds[0][i].Message)
	}
	fmt.Println()

	// Check to make sure all the feeds are the same. This doesn't work since there is no way to guarantee
	// that two posts do not have the same timestamp
	for i := 1; i < numUsers; i++ {
		for k := 0; k < numUsers; k++ {
			if reflect.DeepEqual(feeds[0][k].Message, feeds[i][k].Message) == false {
				if feeds[0][k].Timestamp.Equal(feeds[i][k].Timestamp) == false {
					t.Error(fmt.Sprintf("User %d feed not equal to first feed. Expected %s, Found %s", i, feeds[0][k].Message, feeds[i][k].Message))
				}
			}
		}
	}
}
