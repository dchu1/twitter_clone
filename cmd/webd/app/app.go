package app

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// User struct contatining attributes of a User
type User struct {
	mu        sync.Mutex // protects User
	FirstName string
	LastName  string
	Email     string
	Password  string
<<<<<<< HEAD
	Id        uint64
	following []*User
	followers []*User
=======
	id        uint64
	following map[uint64]uint64
	followers map[uint64]uint64
>>>>>>> master
	post      []*Post
}

// Post struct containing attributes of a Post
type Post struct {
	mu        sync.Mutex // protects Post
	Id        uint64     // This is a unique id. Type might be different depending on how we generate unique ids.
	Timestamp time.Time  // time this post was made
	Message   string     // the text of the post
	UserID    uint64     //id of the user who wrote the post
}

// App struct containig master list of users and posts
type App struct {
	usersMu sync.Mutex // protects users map
	postsMu sync.Mutex // protects posts map
	users   map[uint64]*User
	posts   map[uint64]*Post
	userID  uint64
	postID  uint64
}

// ByTime implements sort.Interface for []Person based on
// the Age field.
type ByTime []*Post

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }

func (u *User) String() string {
	return fmt.Sprintf("FirstName: %s, LastName: %s, Email: %s, Password: %s, id: %d, following: %d, followers: %d, posts: %d",
		u.FirstName, u.LastName, u.Email, u.Password, u.Id, len(u.following), len(u.followers), len(u.post))
}

func MakeApp() *App {
	return &App{sync.Mutex{}, sync.Mutex{}, make(map[uint64]*User), make(map[uint64]*Post), 0, 0}
}

func (appList *App) FollowUser(followingUserID uint64, UserIDToFollow uint64) {
	appList.usersMu.Lock()
	defer appList.usersMu.Unlock()

	appList.users[followingUserID].mu.Lock()
	//Add userID to be followed in the following list of user who wants to follow
	followingUserIDObject := appList.users[followingUserID]
	newfollowing := followingUserIDObject.following
	newfollowing[UserIDToFollow] = UserIDToFollow
	followingUserIDObject.following = newfollowing
	appList.users[followingUserID].mu.Unlock()

	appList.users[UserIDToFollow].mu.Lock()
	//Add userID who is following in the followers list of the user being followed
	UserIDToFollowObject := appList.users[UserIDToFollow]
	newfollowers := UserIDToFollowObject.followers
	newfollowers[followingUserID] = followingUserID
	UserIDToFollowObject.followers = newfollowers
	appList.users[UserIDToFollow].mu.Unlock()
}

func (appList *App) UnFollowUser(followingUserID uint64, UserIDToUnfollow uint64) {
	appList.usersMu.Lock()
	defer appList.usersMu.Unlock()

	appList.users[followingUserID].mu.Lock()
	//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	followingUserIDObject := appList.users[followingUserID]
	newfollowing := followingUserIDObject.following
	delete(newfollowing, UserIDToUnfollow)
	followingUserIDObject.following = newfollowing
	appList.users[followingUserID].mu.Unlock()

	appList.users[UserIDToUnfollow].mu.Lock()
	//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	UserIDToUnfollowObject := appList.users[UserIDToUnfollow]
	newfollowers := UserIDToUnfollowObject.followers
	delete(newfollowers, followingUserID)
	UserIDToUnfollowObject.followers = newfollowers
	appList.users[UserIDToUnfollow].mu.Unlock()
}

func (appList *App) CreatePost(userID uint64, message string) error {
	currTime := time.Now()
	newPost := &Post{sync.Mutex{}, appList.postID, currTime, message, userID}

	appList.postsMu.Lock()
	defer appList.postsMu.Unlock()
	appList.posts[appList.postID] = newPost
	appList.postID++

	// Temporary code
	appList.usersMu.Lock()
	defer appList.usersMu.Unlock()
	appList.users[userID].post = append(appList.users[userID].post, newPost)

	return nil
}

func MakeUser(firstname string, lastname string, email string, password string, id uint64) *User {
	return &User{sync.Mutex{}, firstname, lastname, email, password, id, make(map[uint64]uint64), make(map[uint64]uint64), make([]*Post, 0, 10)}
}

func (appList *App) AddUser(firstname string, lastname string, email string, password string) (uint64, error) {
	appList.usersMu.Lock()
	defer appList.usersMu.Unlock()

	newUser := MakeUser(firstname, lastname, email, password, appList.userID)
	appList.users[appList.userID] = newUser
	appList.userID++
	return newUser.Id, nil
}

func (appList *App) GetUsers() map[uint64]*User {
	return appList.users
}

func (appList *App) GetUser(id uint64) *User {
	return appList.users[id]
}

func (appList *App) GetFeed(userId uint64) ([]*Post, error) {
	// naive implementation
	posts := make([]*Post, 0, 100)
	posts = append(posts, appList.users[userId].post...)
	for _, v := range appList.users[userId].following {
		posts = append(posts, appList.users[v].post...)
	}
	// sort
	sort.Sort(ByTime(posts))
	return posts, nil
}

// Function to strip out data for export, i.e. password, etc...
func CleanUserData(u *User) User {
	ret := *u
	ret.Password = 0
	return ret
}
