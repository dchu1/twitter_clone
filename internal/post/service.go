package post

import (
	"context"
	"sort"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/golang/protobuf/ptypes"
)

type postServiceServer struct {
	postRepo PostRepository
	pb.UnimplementedPostServiceServer
}

// ByTime implements sort.Interface for []Post based on
// the timestamp field.
type ByTime []*pb.Post

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	t1, _ := ptypes.Timestamp(a[i].Timestamp)
	t2, _ := ptypes.Timestamp(a[j].Timestamp)
	return !t1.Before(t2)
}

func (s *postServiceServer) CreatePost(ctx context.Context, p *pb.Post) (*pb.PostID, error) {
	id, err := s.postRepo.CreatePost(ctx, p)
	return &pb.PostID{PostID: id}, err
}
func (s *postServiceServer) GetPost(ctx context.Context, postID *pb.PostID) (*pb.Post, error) {
	return s.postRepo.GetPost(ctx, postID.GetPostID())
}
func (s *postServiceServer) GetPosts(ctx context.Context, postIDs *pb.PostIDs) (*pb.Posts, error) {
	posts, err := s.postRepo.GetPosts(ctx, postIDs.GetPostIDs())
	sort.Sort(ByTime(posts))
	return &pb.Posts{Posts: posts}, err
}
func (s *postServiceServer) GetPostsByAuthors(ctx context.Context, userIDs *pb.UserIDs) (*pb.Posts, error) {
	posts, err := s.postRepo.GetPostsByAuthor(ctx, userIDs.GetUserIDs())
	sort.Sort(ByTime(posts))
	return &pb.Posts{Posts: posts}, err
}

// GetPostServiceServer returns a grpc Server for the post service using the provided PostRepository
func GetPostServiceServer(pr *PostRepository) *postServiceServer {
	return &postServiceServer{postRepo: *pr}
}
