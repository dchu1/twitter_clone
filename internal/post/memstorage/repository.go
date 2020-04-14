package memstorage

import (
	"context"
	"errors"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/golang/protobuf/ptypes"
)

type postRepository struct {
	storage *postStorage
}

func GetPostRepository() post.PostRepository {
	return &postRepository{PostStorage}
}

func NewPostRepository(storage *postStorage) post.PostRepository {
	return &postRepository{storage}
}

// CreatePost inserts a post into our post map
func (postRepo *postRepository) CreatePost(ctx context.Context, p *pb.Post) (uint64, error) {
	postRepo.storage.postsRWMu.Lock()
	defer postRepo.storage.postsRWMu.Unlock()
	postEntry := new(postEntry)
	p.PostID = postRepo.storage.generatePostId()
	postEntry.post = p
	postEntry.post.Timestamp, _ = ptypes.TimestampProto(time.Now())
	postRepo.storage.posts[p.PostID] = postEntry
	return p.PostID, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	postRepo.storage.postsRWMu.RLock()
	defer postRepo.storage.postsRWMu.RUnlock()
	postEntry, exists := postRepo.storage.posts[postID]
	if !exists {
		return nil, errors.New("user not found")
	}
	p := *postEntry.post
	return &p, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPosts(ctx context.Context, postIDs []uint64) ([]*pb.Post, error) {
	postRepo.storage.postsRWMu.RLock()
	defer postRepo.storage.postsRWMu.RUnlock()
	postArr := make([]*pb.Post, 0, len(postIDs))
	for _, v := range postIDs {
		postEntry, exists := postRepo.storage.posts[v]
		if !exists {
			return postArr, errors.New("post not found")
		}
		postArr = append(postArr, postEntry.post)
	}
	return postArr, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPostsByAuthor(ctx context.Context, userIDs []uint64) ([]*pb.Post, error) {
	postRepo.storage.postsRWMu.RLock()
	defer postRepo.storage.postsRWMu.RUnlock()
	postArr := make([]*pb.Post, 0, len(userIDs)*100)
	for _, v := range postRepo.storage.posts {
		v.mu.RLock()
		for _, u := range userIDs {
			if v.post.UserId == u {
				postArr = append(postArr, v.post)
				break
			}
		}
		v.mu.RUnlock()
	}
	return postArr, nil
}

func (postRepo *postRepository) UpdatePost(ctx context.Context, p pb.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(ctx context.Context, postIDs uint64) error {
	return errors.New("Feature not implemented")
}
