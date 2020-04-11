package auth

import (
	"context"
	"errors"
)

type Credentials struct {
	Email    string
	Password string
}

type CredentialsRepository interface {
	CreateCredentials(context.Context, Credentials) error
	GetCredentials(context.Context, Credentials) error
	UpdateCredentials(context.Context, Credentials) error
	DeleteCredentials(context.Context, Credentials) error
}

// Service is the interface that provides auth methods.
type Service interface {
	CreateCredentials(context.Context, Credentials) error
	GetCredentials(context.Context, Credentials) error
	UpdateCredentials(context.Context, Credentials) error
	DeleteCredentials(context.Context, Credentials) error
}

type service struct {
	credentialsRepo CredentialsRepository
	sessionManager  *Manager
}

func NewService(cr CredentialsRepository, sm *Manager) Service {
	return &service{cr, sm}
}

func (s *service) CreateCredentials(ctx context.Context, creds Credentials) error {
	return s.credentialsRepo.CreateCredentials(ctx, creds)
}
func (s *service) GetCredentials(ctx context.Context, creds Credentials) error {
	return s.credentialsRepo.GetCredentials(ctx, creds)
}
func (s *service) UpdateCredentials(ctx context.Context, creds Credentials) error {
	return errors.New("Feature not implemented")
}
func (s *service) DeleteCredentials(ctx context.Context, creds Credentials) error {
	return errors.New("Feature not implemented")
}
