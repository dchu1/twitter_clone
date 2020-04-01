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

// UsersMapCopy creates a deep copy of the users map.
func (appList *App) UsersMapCopy() map[uint64]User {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	cp := make(map[uint64]User)
	for k, v := range appList.users {
		cp[k] = v.Clone()
	}
	return cp
}

// PostsMapCopy creates a deep copy of the posts map
func (appList *App) PostsMapCopy() map[uint64]Post {
	appList.postsRWMu.RLock()
	defer appList.postsRWMu.RUnlock()
	cp := make(map[uint64]Post)
	for k, v := range appList.posts {
		cp[k] = Post{sync.Mutex{}, v.Id, v.Timestamp, v.Message, v.UserID}
	}
	return cp
}

// copyFollowMap makes a deep copy of a user's following or followed map
func copyFollowMap(m map[uint64]uint64) map[uint64]uint64 {
	cp := make(map[uint64]uint64)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// Clone creates a deep copy of a User object
func (u *User) Clone() User {
	u.followingRWMu.RLock()
	u.followersRWMu.RLock()
	defer u.followingRWMu.RUnlock()
	defer u.followersRWMu.RUnlock()
	retUser := MakeUser(u.FirstName, u.LastName, u.Email, u.Id)
	retUser.following = copyFollowMap(u.following)
	retUser.followers = copyFollowMap(u.followers)
	return *retUser
}

// String is the toString() method for a user
func (u *User) String() string {
	return fmt.Sprintf("FirstName: %s, LastName: %s, Email: %s, id: %d, following: %d, followers: %d, posts: %d",
		u.FirstName, u.LastName, u.Email, u.Id, len(u.following), len(u.followers), len(u.post))
}

// MakeApp creates an instsance of an App
func MakeApp() *App {
	return &App{sync.RWMutex{}, sync.RWMutex{}, sync.Mutex{}, sync.Mutex{}, sync.RWMutex{}, make(map[string]string), make(map[uint64]*User), make(map[uint64]*Post), 0, 0}
}

// generateUserId gets a value from the userId counter, then increments the counter
func (appList *App) generateUserId() uint64 {
	appList.userIDMu.Lock()
	defer appList.userIDMu.Unlock()
	uid := appList.userID
	appList.userID++
	return uid
}

// generatePostId gets a value from the postId counter, then increments the counter
func (appList *App) generatePostId() uint64 {
	appList.postIDMu.Lock()
	defer appList.postIDMu.Unlock()
	uid := appList.postID
	appList.postID++
	return uid
}

// FollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user is following another user
func (appList *App) FollowUser(followingUserID uint64, UserIDToFollow uint64) error {

	if followingUserID == UserIDToFollow {
		return errors.New("duplicate user ids")
	}

	//Add userID to be followed in the following list of user who wants to follow
	followingUserIDObject := appList.getUser(followingUserID)
	followingUserIDObject.followingRWMu.Lock()
	newfollowing := followingUserIDObject.following
	newfollowing[UserIDToFollow] = UserIDToFollow
	followingUserIDObject.following = newfollowing
	followingUserIDObject.followingRWMu.Unlock()

	//Add userID who is following in the followers list of the user being followed
	UserIDToFollowObject := appList.getUser(UserIDToFollow)
	UserIDToFollowObject.followersRWMu.Lock()
	newfollowers := UserIDToFollowObject.followers
	newfollowers[followingUserID] = followingUserID
	UserIDToFollowObject.followers = newfollowers
	UserIDToFollowObject.followersRWMu.Unlock()

	return nil
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (appList *App) UnFollowUser(followingUserID uint64, UserIDToUnfollow uint64) error {
	if followingUserID == UserIDToUnfollow {
		return errors.New("duplicate user ids")
	}

	//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	followingUserIDObject := appList.getUser(followingUserID)
	followingUserIDObject.followingRWMu.Lock()
	newfollowing := followingUserIDObject.following
	delete(newfollowing, UserIDToUnfollow)
	followingUserIDObject.following = newfollowing
	followingUserIDObject.followingRWMu.Unlock()

	//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	UserIDToUnfollowObject := appList.getUser(UserIDToUnfollow)
	UserIDToUnfollowObject.followersRWMu.Lock()
	newfollowers := UserIDToUnfollowObject.followers
	delete(newfollowers, followingUserID)
	UserIDToUnfollowObject.followers = newfollowers
	UserIDToUnfollowObject.followersRWMu.Unlock()

	return nil
}

// CreatePost creates a post object and appends it to the appropriate
// data structures
func (appList *App) CreatePost(userID uint64, message string) error {
	currTime := time.Now()
	postId := appList.generatePostId()
	newPost := &Post{sync.Mutex{}, postId, currTime, message, userID}

	appList.postsRWMu.Lock()
	appList.posts[postId] = newPost
	appList.postsRWMu.Unlock()

	user := appList.getUser(userID)
	user.postsRWMu.Lock()
	appList.users[userID].post = append(appList.users[userID].post, newPost)
	user.postsRWMu.Unlock()

	return nil
}

// MakeUser returns a user object
func MakeUser(firstname string, lastname string, email string, id uint64) *User {
	return &User{sync.RWMutex{}, sync.RWMutex{}, sync.RWMutex{}, firstname, lastname, email, id, make(map[uint64]uint64), make(map[uint64]uint64), make([]*Post, 0, 100)}
}

// AddUser adds a user to the appropriate data structures
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
	appList.users[userId] = newUser
	appList.usersRWMu.Unlock()

	appList.credentialsRWMu.Lock()
	defer appList.credentialsRWMu.Unlock()
	appList.credentials[email] = password

	return newUser.Id, nil
}

// GetUsers returns a list of all users
func (appList *App) GetUsers() map[uint64]*User {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	newMap := make(map[uint64]*User)
	for k, v := range appList.users {
		clone := v.Clone()
		newMap[k] = &clone
	}
	return newMap
}

func (appList *App) getUser(id uint64) *User {
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	return appList.users[id]
}

func (appList *App) GetUser(id uint64) *User {
	user := appList.getUser(id)
	temp := user.Clone()
	return &temp
}

// GetUserPosts returns a list of all the posts by a user.
func (appList *App) GetUserPosts(userId uint64) []Post {
	posts := make([]Post, 0, 100)
	user := appList.getUser(userId)

	user.postsRWMu.RLock()
	for _, v := range user.post {
		posts = append(posts, *v)
	}
	user.postsRWMu.RUnlock()
	return posts
}

// ValidateCredentials checks whether the given username and password match
// those stored in the credentials map
func (appList *App) ValidateCredentials(username string, password string) bool {
	appList.credentialsRWMu.RLock()
	defer appList.credentialsRWMu.RUnlock()
	return appList.credentials[username] == password
}

// GetUserByUsername returns a user object by their username
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

// GetFeed returns an array of Posts that represent a given user's feed
func (appList *App) GetFeed(userId uint64) ([]Post, error) {
	posts := make([]Post, 0, 100)

	// Get users posts
	userPosts := appList.GetUserPosts(userId)
	posts = append(posts, userPosts...)

	user := appList.getUser(userId)

	// Get user following posts
	user.followingRWMu.RLock()
	for _, v := range user.following {
		userPosts := appList.GetUserPosts(v)
		posts = append(posts, userPosts...)
	}
	user.followingRWMu.RUnlock()
	// sort
	sort.Sort(ByTime(posts))
	return posts, nil
}

// GetFollowing returns an array of users that the given user is following
func (appList *App) GetFollowing(userId uint64) ([]User, error) {
	// Get the user object from the users map
	user := appList.getUser(userId)

	user.followingRWMu.RLock()
	defer user.followingRWMu.RUnlock()

	tempArray := make([]User, 0, 100)
	for followerId := range user.following {
		tempArray = append(tempArray, appList.getUser(followerId).Clone())
	}
	return tempArray, nil
}

// GetNotFollowing returns an array of users that the given user is not following
func (appList *App) GetNotFollowing(userId uint64) ([]User, error) {
	// Get the user object from the users map
	appList.usersRWMu.RLock()
	defer appList.usersRWMu.RUnlock()
	user := appList.users[userId]

	user.followingRWMu.RLock()
	defer user.followingRWMu.RUnlock()

	tempArray := make([]User, 0, 100)
	for k, v := range appList.users {
		// check if user k exists in the user's following list. If not, add it to our
		// temp array
		_, exists := user.following[k]
		if !exists && k != userId {
			tempArray = append(tempArray, v.Clone())
		}
	}
	return tempArray, nil
}
