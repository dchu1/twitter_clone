package post

import (
	"context"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
)

// Post struct containing attributes of a Post
// type Post struct {
// 	PostID    uint64    // This is a unique id. Type might be different depending on how we generate unique ids.
// 	Timestamp time.Time // time this post was made
// 	Message   string    // the text of the post
// 	UserID    uint64    //id of the user who wrote the post
// }

type PostRepository interface {
	CreatePost(context.Context, *pb.Post) (uint64, error)
	GetPost(context.Context, uint64) (*pb.Post, error)
	GetPosts(context.Context, []uint64) ([]*pb.Post, error)
	GetPostsByAuthor(context.Context, []uint64) ([]*pb.Post, error)
	UpdatePost(context.Context, pb.Post) error
	DeletePost(context.Context, uint64) error
}
