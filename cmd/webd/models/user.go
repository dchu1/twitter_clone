package models

import "time"

// User struct contatining attributes of a User
type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	id        uint64
	following []*User
	followers []*User
	post      []*Post
}

// Post struct containing attributes of a Post
type Post struct {
	id        uint64    // This is a unique id. Type might be different depending on how we generate unique ids.
	timestamp time.Time // time this post was made
	message   string    // the text of the post
	userID    uint64    //id of the user who wrote the post
}

// App struct containig master list of users and posts
type App struct {
	users  map[uint64]*User
	posts  map[uint64]*Post
	userID uint64
	postID uint64
}
