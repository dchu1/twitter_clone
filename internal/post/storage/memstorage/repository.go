package memstorage

import (
	"context"
	"errors"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
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
	result := make(chan uint64, 1)
	errorchan := make(chan error, 1)
	go postRepo.storage.createPost(p, result, errorchan)

	select {
	case postID := <-result:
		return postID, nil
	case err := <-errorchan:
		//Sending 0 as an invalid postID
		return 0, err
	case <-ctx.Done():
		go func() {
			select {
			case <-result:
				postRepo.storage.postsRWMu.Lock()
				delete(postRepo.storage.posts, p.PostID)
				postRepo.storage.postsRWMu.Unlock()
				return
			case <-errorchan:
				return
			}
		}()

		return 0, ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	result := make(chan *pb.Post, 1)
	errorchan := make(chan error, 1)

	go postRepo.storage.getPost(postID, result, errorchan)

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
	result := make(chan []*pb.Post, 1)
	errorchan := make(chan error, 1)

	go postRepo.storage.getPosts(postIDs, result, errorchan)

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
	result := make(chan []*pb.Post, 1)
	errorchan := make(chan error, 1)

	go postRepo.storage.getPostsByAuthor(userIDs, result, errorchan)

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

func (postRepo *postRepository) DeletePost(ctx context.Context, postID uint64) error {
	errorchan := make(chan error, 1)
	buffer := make(chan *postEntry, 1)
	go postRepo.storage.deletePost(postID, errorchan, buffer)

	select {
	case err := <-errorchan:
		return err
	case <-ctx.Done():
		// if ctx done, need to continue to listen to know whether to add postEntry back into db
		go func() {
			select {
			case err := <-errorchan:
				// if result != nil, an error occurred and so don't need to add back into db
				if err != nil {
					return
				}
				postEntry := <-buffer
				postRepo.storage.postsRWMu.Lock()
				defer postRepo.storage.postsRWMu.Unlock()
				postRepo.storage.posts[postID] = postEntry
				return
			}

		}()
		return ctx.Err()
	}
}
