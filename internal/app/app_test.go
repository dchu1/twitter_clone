package app

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/storage/memstorage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
)

func TestAddUser(t *testing.T) {
	storage := memstorage.NewMemoryStorage()
	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
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
	storage := memstorage.NewMemoryStorage()
	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
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
	storage := memstorage.NewMemoryStorage()
	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
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

func TestCreatePost(t *testing.T) {
	storage := memstorage.NewMemoryStorage()
	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
	userID, _ := app.CreateUser(nil, user.AccountInformation{"FollowerFirst", "FollowerLast", "follower@test.com", 0})
	postID, _ := app.CreatePost(nil, user.AccountInformation{"FollowerFirst", "FollowerLast", "follower@test.com", userID}, post.Post{Message: "Test Message"})

	posts, err := app.GetFeed(nil, userID)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(posts[0].Message, "Test Message") || !reflect.DeepEqual(posts[0].PostID, postID) {
		t.Error("Test Failed Posts struct not in sync")
	}
}

// func TestConcurrentGetUsers(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var userlist []map[uint64]*User
// 	userListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage := memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", 0)

// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			userListmu.Lock()
// 			defer userListmu.Unlock()
// 			userlist = append(userlist, app.GetUsers())
// 		}(user)
// 	}
// 	wg.Wait()
// 	for _, user := range userlist {
// 		if reflect.DeepEqual(app.users, user) == false {
// 			t.Error("Error getting users")
// 		}
// 	}

// }

// func TestConcurrentGetUser(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var userlist []*User
// 	userListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", "TestPassword")
// 	userID := app.users[0].Id
// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			userListmu.Lock()
// 			defer userListmu.Unlock()
// 			userlist = append(userlist, app.getUser(userID))
// 		}(user)
// 	}
// 	wg.Wait()
// 	for _, user := range userlist {
// 		if reflect.DeepEqual(app.users[0], user) == false {
// 			t.Error("Error getting user")
// 		}
// 	}

// }

// func TestConcurrentGetUserPosts(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var postlist [][]Post
// 	postlistmu := sync.Mutex{}
// 	numPost := 100
// 	wg.Add(numPost)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", "TestPassword")
// 	userID := app.users[0].Id
// 	app.CreatePost(userID, "TestMessage")
// 	for post := 0; post < numPost; post++ {
// 		go func(post int) {
// 			defer wg.Done()
// 			postlistmu.Lock()
// 			defer postlistmu.Unlock()
// 			postlist = append(postlist, app.GetUserPosts(userID))
// 		}(post)
// 	}
// 	wg.Wait()
// 	for _, post := range postlist {
// 		if post[0].Message != "TestMessage" {
// 			t.Error("Incorrect message for the post")
// 		}
// 		if post[0].UserID != userID {
// 			t.Error("Incorrect userID for the post")
// 		}
// 	}

// }

// func TestConcurrentValidateCredentials(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var resultList []bool
// 	resultListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", "TestPassword")
// 	username := app.users[0].Email
// 	password := app.credentials[username]

// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			resultListmu.Lock()
// 			defer resultListmu.Unlock()
// 			resultList = append(resultList, app.ValidateCredentials(username, password))
// 		}(user)
// 	}
// 	wg.Wait()
// 	for _, result := range resultList {
// 		if result != true {
// 			t.Error("Error validating user credentials")
// 		}
// 	}

// }

// func TestConcurrentGetUserByUsername(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var userList []*User
// 	userListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", "TestPassword")
// 	username := app.users[0].Email

// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			userObject, _ := app.GetUserByUsername(username)
// 			userListmu.Lock()
// 			defer userListmu.Unlock()
// 			userList = append(userList, userObject)
// 		}(user)
// 	}
// 	wg.Wait()
// 	for _, user := range userList {
// 		if reflect.DeepEqual(app.users[0], user) != true {
// 			t.Error("Error retreiving user given the email")
// 		}
// 	}

// }

// func TestConcurrentGetFollowing(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var userList [][]User
// 	userListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName1", "TestLastName1", "TestEmail1", "TestPassword1")
// 	app.AddUser("TestFirstName2", "TestLastName2", "TestEmail2", "TestPassword2")

// 	followingUserID := app.users[0].Id
// 	UserIDToFollow := app.users[1].Id

// 	app.FollowUser(followingUserID, UserIDToFollow)

// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			userObject, _ := app.GetFollowing(followingUserID)
// 			userListmu.Lock()
// 			defer userListmu.Unlock()
// 			userList = append(userList, userObject)
// 		}(user)
// 	}
// 	wg.Wait()
// 	appUserObject := app.users[1]
// 	for _, user := range userList {
// 		if user[0].Id != UserIDToFollow {
// 			t.Error("Incorrect user ID of the user which is being followed")
// 		}
// 		if user[0].FirstName != appUserObject.FirstName {
// 			t.Error("Incorrect First Name of the user which is being followed")
// 		}
// 		if user[0].LastName != appUserObject.LastName {
// 			t.Error("Incorrect Last Name of the user which is being followed")
// 		}
// 		if user[0].Email != appUserObject.Email {
// 			t.Error("Incorrect Email of the user which is being followed")
// 		}
// 	}

// }

// func TestConcurrentGetNotFollowing(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var userList [][]User
// 	userListmu := sync.Mutex{}
// 	numUser := 100
// 	wg.Add(numUser)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName1", "TestLastName1", "TestEmail1", "TestPassword1")
// 	app.AddUser("TestFirstName2", "TestLastName2", "TestEmail2", "TestPassword2")
// 	app.AddUser("TestFirstName3", "TestLastName3", "TestEmail3", "TestPassword3")

// 	followingUserID := app.users[0].Id
// 	UserIDToFollow := app.users[1].Id

// 	app.FollowUser(followingUserID, UserIDToFollow)

// 	for user := 0; user < numUser; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			userObject, _ := app.GetNotFollowing(followingUserID)
// 			userListmu.Lock()
// 			defer userListmu.Unlock()
// 			userList = append(userList, userObject)
// 		}(user)
// 	}
// 	wg.Wait()
// 	appUserObject := app.users[2]
// 	for _, user := range userList {
// 		if user[0].Id != appUserObject.Id {
// 			t.Error("Incorrect user ID of the user which is not being followed")
// 		}
// 		if user[0].FirstName != appUserObject.FirstName {
// 			t.Error("Incorrect First Name of the user which is not being followed")
// 		}
// 		if user[0].LastName != appUserObject.LastName {
// 			t.Error("Incorrect Last Name of the user which is not being followed")
// 		}
// 		if user[0].Email != appUserObject.Email {
// 			t.Error("Incorrect Email of the user which is not being followed")
// 		}
// 	}

// }

// func TestConcurrentGetFeed(t *testing.T) {
// 	var wg sync.WaitGroup
// 	var postList [][]Post
// 	postListmu := sync.Mutex{}
// 	numPost := 100
// 	wg.Add(numPost)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirstName", "TestLastName", "TestEmail", "TestPassword")
// 	userID := app.users[0].Id
// 	app.CreatePost(userID, "TestMessage")

// 	for post := 0; post < numPost; post++ {
// 		go func(post int) {
// 			defer wg.Done()
// 			postObject, _ := app.GetFeed(userID)
// 			postListmu.Lock()
// 			defer postListmu.Unlock()
// 			postList = append(postList, postObject)
// 		}(post)
// 	}
// 	wg.Wait()
// 	for _, post := range postList {
// 		if post[0].Message != "TestMessage" {
// 			t.Error("Retreive feed has incorrect message value")
// 		}
// 		if post[0].UserID != userID {
// 			t.Error("Retreive feed has incorrect userID")
// 		}
// 	}

// }

// func TestConcurrentGenerateUserId(t *testing.T) {
// 	var wg sync.WaitGroup
// 	numUsers := 100
// 	wg.Add(numUsers)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	for i := 0; i < numUsers; i++ {
// 		go func() {
// 			defer wg.Done()
// 			app.generateUserId()
// 		}()
// 	}
// 	wg.Wait()
// 	if app.userID != 100 {
// 		t.Error("user ID value not updated correctly in the App struct")
// 	}
// }

// func TestConcurrentGeneratePostId(t *testing.T) {
// 	var wg sync.WaitGroup
// 	numPosts := 100
// 	wg.Add(numPosts)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	for i := 0; i < numPosts; i++ {
// 		go func() {
// 			defer wg.Done()
// 			app.generatePostId()
// 		}()
// 	}
// 	wg.Wait()
// 	if app.postID != 100 {
// 		t.Error("post ID value not updated correctly in the App struct")
// 	}
// }

// func TestConcurrentCreatePost(t *testing.T) {
// 	var wg sync.WaitGroup
// 	numPosts := 100
// 	wg.Add(numPosts)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	app.AddUser("TestFirst", "TestLast", "Test@test.com", "Testpass")
// 	userID := app.users[0].Id
// 	for post := 0; post < numPosts; post++ {
// 		go func(post int) {
// 			defer wg.Done()
// 			message := "TestMessage " + strconv.Itoa(post)
// 			app.CreatePost(userID, message)
// 		}(post)
// 	}
// 	wg.Wait()
// 	if len(app.posts) != 100 {
// 		t.Error("Not all posts added to the app struct posts map")
// 	}
// 	if app.postID != 100 {
// 		t.Error("post ID value not updated in the App struct")
// 	}
// }

// func TestConcurrentAddUser(t *testing.T) {
// 	var wg sync.WaitGroup
// 	numUsers := 100
// 	wg.Add(numUsers)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	for user := 0; user < numUsers; user++ {
// 		go func(user int) {
// 			defer wg.Done()
// 			firstName := "TestFirstName" + strconv.Itoa(user)
// 			lastName := "TestLastName" + strconv.Itoa(user)
// 			email := "TestEmail" + strconv.Itoa(user)
// 			password := "TestPassword" + strconv.Itoa(user)
// 			app.AddUser(firstName, lastName, email, password)
// 		}(user)
// 	}
// 	wg.Wait()
// 	if len(app.users) != 100 {
// 		t.Error("All users not added in the struct")
// 	}
// }

func TestConcurrentFollow(t *testing.T) {
	var wg sync.WaitGroup
	rand.Seed(42)
	numUsers := 100
	wg.Add(numUsers)
	storage := memstorage.NewMemoryStorage()
	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
	// Create users
	for i := 0; i < numUsers; i++ {
		app.CreateUser(nil, user.AccountInformation{strconv.Itoa(i), strconv.Itoa(i), strconv.Itoa(i), 0})
	}

	// Have them all follow each other and then make a post
	for i := 0; i < numUsers; i++ {
		go func(userId uint64) {
			defer wg.Done()
			for k := 0; k < numUsers; k++ {
				app.FollowUser(nil, uint64(userId), uint64(k))
			}
			// sleep for some random amount of time between 0 and 5 seconds
			//time.Sleep(time.Second * time.Duration(rand.Intn(20)))
			app.CreatePost(nil, user.AccountInformation{UserID: uint64(userId)}, post.Post{Message: strconv.FormatUint(userId, 10)})
		}(uint64(i))
	}

	wg.Wait()

	// Get each users following list
	following := make([][]*user.User, numUsers)
	for i := 0; i < numUsers; i++ {
		var err error
		following[i], err = app.GetFollowing(nil, uint64(i))
		if err != nil {
			t.Error(err.Error())
		}
		if len(following[i]) != numUsers-1 {
			t.Error(fmt.Sprintf("Incorrect number of users following. Expected %d, found %d", numUsers-1, len(following[i])))
		}
	}

	// Get each users feed
	feeds := make([][]*Post, numUsers)
	for i := 0; i < numUsers; i++ {
		var err error
		feeds[i], err = app.GetFeed(nil, uint64(i))
		if err != nil {
			panic(err.Error())
		}
		if len(feeds[i]) != numUsers {
			t.Error(fmt.Sprintf("Not enough posts in user %d feed. Expected %d, found %d", i, numUsers, len(feeds[i])))
		}
	}

	// Check to make sure all the feeds are the same. This doesn't work since there is no way to guarantee
	// that two posts do not have the same timestamp
	for i := 1; i < numUsers; i++ {
		for k := 0; k < numUsers; k++ {
			if reflect.DeepEqual(feeds[0][k].Message, feeds[i][k].Message) == false {
				if feeds[0][k].Timestamp.Equal(feeds[i][k].Timestamp) == false {
					t.Error(fmt.Sprintf("User %d feed not equal to first feed", i))
				}
			}
		}
	}
}

// func TestConcurrentUnfollowUser(t *testing.T) {
// 	var wg sync.WaitGroup
// 	actualFollowing := make(map[uint64]uint64)
// 	numUsers := 100
// 	wg.Add(numUsers)
// 	storage = memstorage.NewMemoryStorage()
// 	app := NewService(memstorage.NewUserRepository(storage), memstorage.NewPostRepository(storage))
// 	// Create users
// 	for user := 0; user < numUsers; user++ {
// 		app.AddUser("TestFirstName"+strconv.Itoa(user), "TestLastName"+strconv.Itoa(user), "TestEmail"+strconv.Itoa(user), "TestPassword"+strconv.Itoa(user))
// 	}

// 	for followingUserID := 0; followingUserID < numUsers; followingUserID++ {
// 		for UserIDToFollow := 0; UserIDToFollow < numUsers; UserIDToFollow++ {
// 			app.FollowUser(uint64(followingUserID), uint64(UserIDToFollow))
// 		}
// 	}

// 	for user := 0; user < numUsers; user++ {
// 		if len(app.users[uint64(user)].following) != 99 {
// 			t.Error("Incorrect following list for user" + strconv.Itoa(user))
// 		}
// 	}

// 	for followingUserID := 0; followingUserID < numUsers; followingUserID++ {
// 		go func(followingUserID int) {
// 			defer wg.Done()
// 			for UserIDToUnfollow := 0; UserIDToUnfollow < numUsers; UserIDToUnfollow++ {
// 				app.UnFollowUser(uint64(followingUserID), uint64(UserIDToUnfollow))
// 			}
// 		}(followingUserID)
// 	}
// 	wg.Wait()
// 	for _, userObject := range app.users {
// 		if reflect.DeepEqual(userObject.following, actualFollowing) == false {
// 			t.Error("Unsuccessful Unfollow operation ")
// 		}
// 	}
// }
