package auth

import "errors"

type Credentials struct {
	Email    string
	Password string
}

type CredentialsRepository interface {
	CreateCredentials(Credentials) error
	GetCredentials(Credentials) error
	UpdateCredentials(Credentials) error
	DeleteCredentials(Credentials) error
}

// Service is the interface that provides auth methods.
type Service interface {
	CreateCredentials(Credentials) error
	GetCredentials(Credentials) error
	UpdateCredentials(Credentials) error
	DeleteCredentials(Credentials) error
}

type service struct {
	credentialsRepo CredentialsRepository
}

func NewService(cr CredentialsRepository) Service {
	return &service{cr}
}

func (s *service) CreateCredentials(creds Credentials) error {
	return s.credentialsRepo.CreateCredentials(creds)
}
func (s *service) GetCredentials(creds Credentials) error {
	return s.credentialsRepo.GetCredentials(creds)
}
func (s *service) UpdateCredentials(creds Credentials) error {
	return errors.New("Feature not implemented")
}
func (s *service) DeleteCredentials(creds Credentials) error {
	return errors.New("Feature not implemented")
}
