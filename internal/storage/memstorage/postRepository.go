package memstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
)

type postRepository struct {
	storage *memoryStorage
}

type postEntry struct {
	mu   sync.Mutex // protects Post
	post *post.Post
}

func NewPostRepository(storage *memoryStorage) post.PostRepository {
	return &postRepository{storage}
}

// CreatePost inserts a post into our post map
func (postRepo *postRepository) CreatePost(post post.Post) (uint64, error) {
	postRepo.storage.postsRWMu.Lock()
	defer postRepo.storage.postsRWMu.Unlock()
	postEntry := new(postEntry)
	post.PostID = postRepo.storage.generatePostId()
	postEntry.post = &post
	postRepo.storage.posts[post.PostID] = postEntry
	fmt.Printf("%v", postRepo.storage.posts)
	return post.PostID, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPosts(postIDs []uint64) ([]*post.Post, error) {
	postRepo.storage.postsRWMu.RLock()
	defer postRepo.storage.postsRWMu.RUnlock()
	postArr := make([]*post.Post, 0, len(postIDs))
	for _, v := range postIDs {
		postArr = append(postArr, postRepo.storage.posts[v].post)
	}
	return postArr, nil
}

func (postRepo *postRepository) UpdatePost(post post.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(postIDs uint64) error {
	return errors.New("Feature not implemented")
}
