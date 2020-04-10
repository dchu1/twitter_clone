package post

import "time"

// Post struct containing attributes of a Post
type Post struct {
	PostID    uint64    // This is a unique id. Type might be different depending on how we generate unique ids.
	Timestamp time.Time // time this post was made
	Message   string    // the text of the post
	UserID    uint64    //id of the user who wrote the post
}

type PostRepository interface {
	CreatePost(Post) (uint64, error)
	GetPosts([]uint64) ([]*Post, error)
	UpdatePost(Post) error
	DeletePost(uint64) error
}
