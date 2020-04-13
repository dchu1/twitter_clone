package post

import (
	"context"
	"errors"
)

// Service is the interface that provides app methods.
type Service interface {
	CreatePost(context.Context, Post) (uint64, error)
	GetPost(context.Context, uint64) (*Post, error)
	GetPosts(context.Context, []uint64) ([]*Post, error)
	UpdatePost(context.Context, Post) error
	DeletePost(context.Context, uint64) error
}

type service struct {
	postRepo PostRepository
}

func NewService(pr PostRepository) Service {
	return &service{pr}
}

func (s *service) CreatePost(ctx context.Context, p Post) (uint64, error) {
	return s.postRepo.CreatePost(ctx, p)
}
func (s *service) GetPost(ctx context.Context, postID uint64) (*Post, error) {
	return s.postRepo.GetPost(ctx, postID)
}
func (s *service) GetPosts(ctx context.Context, postIDs []uint64) ([]*Post, error) {
	return s.postRepo.GetPosts(ctx, postIDs)
}
func (s *service) UpdatePost(ctx context.Context, p Post) error {
	return errors.New("Feature not implemented")
}
func (s *service) DeletePost(ctx context.Context, postID uint64) error {
	return errors.New("Feature not implemented")
}
