package user

import (
	"context"
	"errors"
)

// Service is the interface that provides app methods.
type Service interface {
	CreateUser(context.Context, AccountInformation) (uint64, error)
	GetUser(context.Context, uint64) (*User, error)
	GetUsers(context.Context, []uint64) ([]*User, error)
	GetFollowing(context.Context, uint64) ([]*User, error)
	GetNotFollowing(context.Context, uint64) ([]*User, error)
	UpdateUserAccountInfo(context.Context, AccountInformation) error
	FollowUser(context.Context, uint64, uint64) error
	UnFollowUser(context.Context, uint64, uint64) error
	DeleteUser(context.Context, uint64) error
}

type service struct {
	userRepo UserRepository
}

func NewService(ur UserRepository) Service {
	return &service{ur}
}

func (s *service) CreateUser(ctx context.Context, info AccountInformation) (uint64, error) {
	return s.userRepo.CreateUser(ctx, info)
}
func (s *service) GetUser(ctx context.Context, userID uint64) (*User, error) {
	return s.userRepo.GetUser(ctx, userID)
}
func (s *service) GetUsers(ctx context.Context, userIDs []uint64) ([]*User, error) {
	return s.userRepo.GetUsers(ctx, userIDs)
}
func (s *service) GetFollowing(ctx context.Context, userID uint64) ([]*User, error) {
	return s.userRepo.GetFollowing(ctx, userID)
}
func (s *service) GetNotFollowing(ctx context.Context, userID uint64) ([]*User, error) {
	return s.userRepo.GetNotFollowing(ctx, userID)
}
func (s *service) UpdateUserAccountInfo(ctx context.Context, info AccountInformation) error {
	return errors.New("Feature not implemented")
}
func (s *service) FollowUser(ctx context.Context, source uint64, target uint64) error {
	return s.userRepo.FollowUser(ctx, source, target)
}
func (s *service) UnFollowUser(ctx context.Context, source uint64, target uint64) error {
	return s.userRepo.UnFollowUser(ctx, source, target)
}
func (s *service) DeleteUser(ctx context.Context, userID uint64) error {
	return errors.New("Feature not implemented")
}
