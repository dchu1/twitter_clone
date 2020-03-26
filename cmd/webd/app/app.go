package app

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// User struct contatining attributes of a User
type User struct {
	followingRWMu sync.RWMutex // protects following map
	followersRWMu sync.RWMutex // protects followers map
	postsRWMu     sync.RWMutex // protects posts map
	FirstName     string       `json:"firstname,omitempty"`
	LastName      string       `json:"lastname,omitempty"`
	Email         string       `json:"email,omitempty"`
	Id            uint64       `json:"userId"`
	following     map[uint64]uint64
	followers     map[uint64]uint64
	post          []*Post
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
	usersRWMu       sync.RWMutex // protects users map
	postsRWMu       sync.RWMutex // protects posts map
	userIDMu        sync.Mutex   // protects userID counter
	postIDMu        sync.Mutex   // protects postID counter
	credentialsRWMu sync.RWMutex // protects the credentials map
	credentials     map[string]string
	users           map[uint64]*User
	posts           map[uint64]*Post
	userID          uint64
	postID          uint64
}

// ByTime implements sort.Interface for []Person based on
// the Age field.
type ByTime []Post

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return !a[i].Timestamp.Before(a[j].Timestamp) }

func (a *App) UsersMapCopy() map[uint64]User {
	a.usersRWMu.RLock()
	defer a.usersRWMu.RUnlock()
	cp := make(map[uint64]User)
	for k, v := range a.users {
		cp[k] = v.Clone()
	}
	return cp
}

func (a *App) PostsMapCopy() map[uint64]Post {
	a.postsRWMu.RLock()
	defer a.postsRWMu.RUnlock()
	cp := make(map[uint64]Post)
	for k, v := range a.posts {
		cp[k] = *v
	}
	return cp
}

func copyFollowMap(m map[uint64]uint64) map[uint64]uint64 {
	cp := make(map[uint64]uint64)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func (u *User) Clone() User {
	u.followingRWMu.RLock()
	u.followersRWMu.RLock()
	defer u.followingRWMu.RUnlock()
	defer u.followersRWMu.RUnlock()
	retUser := *u
	retUser.following = copyFollowMap(u.following)
	retUser.followers = copyFollowMap(u.followers)
	return retUser
}

func (u *User) String() string {
	return fmt.Sprintf("FirstName: %s, LastName: %s, Email: %s, id: %d, following: %d, followers: %d, posts: %d",
		u.FirstName, u.LastName, u.Email, u.Id, len(u.following), len(u.followers), len(u.post))
}

func MakeApp() *App {
	return &App{sync.RWMutex{}, sync.RWMutex{}, sync.Mutex{}, sync.Mutex{}, sync.RWMutex{}, make(map[string]string), make(map[uint64]*User), make(map[uint64]*Post), 0, 0}
}

func (appList *App) generateUserId() uint64 {
	appList.userIDMu.Lock()
	defer appList.userIDMu.Unlock()
	uid := appList.userID
	appList.userID++
	return uid
}

func (appList *App) generatePostId() uint64 {
	appList.postIDMu.Lock()
	defer appList.postIDMu.Unlock()
	uid := appList.postID
	appList.postID++
	return uid
}

func (appList *App) FollowUser(followingUserID uint64, UserIDToFollow uint64) {
	// appList.usersRWMu.Lock()
	// defer appList.usersRWMu.Unlock()

	appList.users[followingUserID].followingRWMu.Lock()
	//Add userID to be followed in the following list of user who wants to follow
	followingUserIDObject := appList.users[followingUserID]
	newfollowing := followingUserIDObject.following
	newfollowing[UserIDToFollow] = UserIDToFollow
	followingUserIDObject.following = newfollowing
	appList.users[followingUserID].followingRWMu.Unlock()

	appList.users[UserIDToFollow].followersRWMu.Lock()
	//Add userID who is following in the followers list of the user being followed
	UserIDToFollowObject := appList.users[UserIDToFollow]
	newfollowers := UserIDToFollowObject.followers
	newfollowers[followingUserID] = followingUserID
	UserIDToFollowObject.followers = newfollowers
	appList.users[UserIDToFollow].followersRWMu.Unlock()
}

func (appList *App) UnFollowUser(followingUserID uint64, UserIDToUnfollow uint64) {
	// appList.usersRWMu.Lock()
	// defer appList.usersRWMu.Unlock()

	appList.users[followingUserID].followingRWMu.Lock()
	//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	followingUserIDObject := appList.users[followingUserID]
	newfollowing := followingUserIDObject.following
	delete(newfollowing, UserIDToUnfollow)
	followingUserIDObject.following = newfollowing
	appList.users[followingUserID].followingRWMu.Unlock()

	appList.users[UserIDToUnfollow].followersRWMu.Lock()
	//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	UserIDToUnfollowObject := appList.users[UserIDToUnfollow]
	newfollowers := UserIDToUnfollowObject.followers
	delete(newfollowers, followingUserID)
	UserIDToUnfollowObject.followers = newfollowers
	appList.users[UserIDToUnfollow].followersRWMu.Unlock()
}

func (appList *App) CreatePost(userID uint64, message string) error {
	currTime := time.Now()
	postId := appList.generatePostId()
	newPost := &Post{sync.Mutex{}, postId, currTime, message, userID}

	appList.postsRWMu.Lock()
	appList.posts[postId] = newPost
	appList.postsRWMu.Unlock()

	// Temporary code
	appList.usersRWMu.Lock()
	defer appList.usersRWMu.Unlock()
	appList.users[userID].post = append(appList.users[userID].post, newPost)

	return nil
}

func MakeUser(firstname string, lastname string, email string, id uint64) *User {
	return &User{sync.RWMutex{}, sync.RWMutex{}, sync.RWMutex{}, firstname, lastname, email, id, make(map[uint64]uint64), make(map[uint64]uint64), make([]*Post, 0, 100)}
}

func (appList *App) AddUser(firstname string, lastname string, email string, password string) (uint64, error) {
	// Check whether user already exists
	user, _ := appList.GetUserByUsername(email)
	if user != nil {
		// cannot return a nil value, so need to make sure that whoever calls this function checks the error, and
		// doesn't just use the value returned!
		return 0, errors.New("duplicate email")
	}

	userId := appList.generateUserId()
	newUser := MakeUser(firstname, lastname, email, userId)

	appList.usersRWMu.Lock()
	defer appList.usersRWMu.Unlock()
	appList.users[userId] = newUser

	appList.credentialsRWMu.Lock()
	defer appList.credentialsRWMu.Unlock()
	appList.credentials[email] = password

	return newUser.Id, nil
}

func (appList *App) GetUsers() map[uint64]*User {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	return appList.users
}

func (appList *App) GetUser(id uint64) *User {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	return appList.users[id]
}

func (appList *App) ValidateCredentials(username string, password string) bool {
	appList.credentialsRWMu.RLock()
	defer appList.credentialsRWMu.RUnlock()
	return appList.credentials[username] == password
}

func (appList *App) GetUserByUsername(email string) (*User, error) {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	for _, v := range appList.users {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, errors.New("user not found")
}

func (appList *App) GetFeed(userId uint64) ([]Post, error) {
	// naive implementation
	posts := make([]Post, 0, 100)
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()

	appList.users[userId].postsRWMu.RLock()
	//posts = append(posts, appList.users[userId].post...)
	for _, v := range appList.users[userId].post {
		posts = append(posts, *v)
	}

	appList.users[userId].postsRWMu.RUnlock()

	for _, v := range appList.users[userId].following {
		appList.users[v].postsRWMu.RLock()
		//posts = append(posts, appList.users[v].post...)
		for _, v := range appList.users[userId].post {
			posts = append(posts, *v)
		}
		appList.users[v].postsRWMu.RUnlock()
	}
	// sort
	sort.Sort(ByTime(posts))
	return posts, nil
}

func (appList *App) GetFollowing(userId uint64) ([]User, error) {
	// Get the user object from the users map
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	user := appList.users[userId]

	tempArray := make([]User, 0, 100)
	for k := range user.following {
		tempArray = append(tempArray, appList.users[k].Clone())
	}
	return tempArray, nil
}
