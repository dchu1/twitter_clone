package memstorage

import (
	"context"
	"errors"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
)

type testPostRepository struct {
	postRepo  post.PostRepository
	sleeptime time.Duration
}

const sleeptime = "5s"

// GetUserRepository returns a UserRepository that uses package level storage
func GetTestPostRepository() post.PostRepository {
	pr := GetPostRepository()
	st, _ := time.ParseDuration(sleeptime)
	return &testPostRepository{pr, st}
}

// NewUserRepository reutnrs a UserRepository that uses the given storage
func NewTestPostRepository(pr post.PostRepository) post.PostRepository {
	st, _ := time.ParseDuration(sleeptime)
	return &testPostRepository{pr, st}
}

// CreatePost inserts a post into our post map
func (testRepo *testPostRepository) CreatePost(ctx context.Context, p *pb.Post) (uint64, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.postRepo.CreatePost(ctx, p)
}

// GetPosts retrieves an array of post from the post map
func (testRepo *testPostRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.postRepo.GetPost(ctx, postID)
}

// GetPosts retrieves an array of post from the post map
func (testRepo *testPostRepository) GetPosts(ctx context.Context, postIDs []uint64) ([]*pb.Post, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.postRepo.GetPosts(ctx, postIDs)
}

// GetPosts retrieves an array of post from the post map
func (testRepo *testPostRepository) GetPostsByAuthor(ctx context.Context, userIDs []uint64) ([]*pb.Post, error) {
	time.Sleep(testRepo.sleeptime)
	return testRepo.postRepo.GetPostsByAuthor(ctx, userIDs)
}

func (testRepo *testPostRepository) UpdatePost(ctx context.Context, p pb.Post) error {
	return errors.New("Feature not implemented")
}

func (testRepo *testPostRepository) DeletePost(ctx context.Context, postID uint64) error {
	time.Sleep(testRepo.sleeptime)
	return testRepo.postRepo.DeletePost(ctx, postID)
}
