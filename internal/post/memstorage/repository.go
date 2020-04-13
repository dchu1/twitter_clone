package memstorage

import (
	"context"
	"errors"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
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
func (postRepo *postRepository) CreatePost(ctx context.Context, post post.Post) (uint64, error) {
	postRepo.storage.postsRWMu.Lock()
	defer postRepo.storage.postsRWMu.Unlock()
	postEntry := new(postEntry)
	post.PostID = postRepo.storage.generatePostId()
	postEntry.post = &post
	postEntry.post.Timestamp = time.Now()
	postRepo.storage.posts[post.PostID] = postEntry
	return post.PostID, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*post.Post, error) {
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
func (postRepo *postRepository) GetPosts(ctx context.Context, postIDs []uint64) ([]*post.Post, error) {
	postRepo.storage.postsRWMu.RLock()
	defer postRepo.storage.postsRWMu.RUnlock()
	postArr := make([]*post.Post, 0, len(postIDs))
	for _, v := range postIDs {
		postArr = append(postArr, postRepo.storage.posts[v].post)
	}
	return postArr, nil
}

func (postRepo *postRepository) UpdatePost(ctx context.Context, post post.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(ctx context.Context, postIDs uint64) error {
	return errors.New("Feature not implemented")
}
