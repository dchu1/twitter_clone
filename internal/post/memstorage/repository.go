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
	result := make(chan uint64)
	errorchan := make(chan error)
	go func() {

		postRepo.storage.postsRWMu.Lock()
		defer postRepo.storage.postsRWMu.Unlock()
		postEntry := new(postEntry)
		p.PostID = postRepo.storage.generatePostId()
		postEntry.post = p
		postEntry.post.Timestamp, _ = ptypes.TimestampProto(time.Now())
		postRepo.storage.posts[p.PostID] = postEntry
		result <- p.PostID
		errorchan <- nil

	}()

	select {
	case postID := <-result:
		return postID, nil
	case err := <-errorchan:
		//Sending 0 as an invalid postID
		return 0, err
	case <-ctx.Done():
		delete(postRepo.storage.posts, p.PostID)
		return 0, ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	result := make(chan *pb.Post)
	errorchan := make(chan error)

	go func() {
		postRepo.storage.postsRWMu.RLock()
		defer postRepo.storage.postsRWMu.RUnlock()
		postEntry, exists := postRepo.storage.posts[postID]
		if !exists {
			result <- nil
			errorchan <- errors.New("user not found")
		} else {
			p := *postEntry.post
			result <- &p
			errorchan <- nil
		}

	}()

	select {
	case post := <-result:
		return post, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPosts(ctx context.Context, postIDs []uint64) ([]*pb.Post, error) {
	result := make(chan []*pb.Post)
	errorchan := make(chan error)

	go func() {

		postRepo.storage.postsRWMu.RLock()
		defer postRepo.storage.postsRWMu.RUnlock()
		postArr := make([]*pb.Post, 0, len(postIDs))
		for _, v := range postIDs {
			postEntry, _ := postRepo.storage.posts[v]
			postArr = append(postArr, postEntry.post)
		}
		result <- postArr
		errorchan <- nil

	}()

	select {
	case posts := <-result:
		return posts, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPostsByAuthor(ctx context.Context, userIDs []uint64) ([]*pb.Post, error) {
	result := make(chan []*pb.Post)
	errorchan := make(chan error)

	go func() {
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
		result <- postArr
		errorchan <- nil
	}()

	select {
	case posts := <-result:
		return posts, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (postRepo *postRepository) UpdatePost(ctx context.Context, p pb.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(ctx context.Context, postIDs uint64) error {
	return errors.New("Feature not implemented")
}
