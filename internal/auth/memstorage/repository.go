package memstorage

import (
	"context"
	"errors"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth"
)

type credentialsRepository struct {
	storage *credentialsStorage
}

type credentialsEntry struct {
	mu          sync.Mutex // protects credentials
	credentials *auth.Credentials
}

func GetCredentialsRepository() auth.CredentialsRepository {
	return &credentialsRepository{CredentialsStorage}
}

func NewCredentialsRepository(storage *credentialsStorage) auth.CredentialsRepository {
	return &credentialsRepository{storage}
}

func (cr *credentialsRepository) CreateCredentials(ctx context.Context, credentials auth.Credentials) error {
	cr.storage.credentialsRWMu.Lock()
	defer cr.storage.credentialsRWMu.Unlock()
	credEntry := new(credentialsEntry)
	credEntry.credentials = &credentials
	cr.storage.credentials[credentials.Email] = credEntry
	return nil
}
func (cr *credentialsRepository) GetCredentials(ctx context.Context, credentials auth.Credentials) (auth.Credentials, error) {
	cr.storage.credentialsRWMu.Lock()
	defer cr.storage.credentialsRWMu.Unlock()
	entry, exists := cr.storage.credentials[credentials.Email]
	if !exists {
		return auth.Credentials{}, errors.New("username not found")
	}
	return auth.Credentials{entry.credentials.Email, entry.credentials.Password}, nil
}

func (cr *credentialsRepository) UpdateCredentials(ctx context.Context, credentials auth.Credentials) error {
	return errors.New("Feature not implemented")
}
func (cr *credentialsRepository) DeleteCredentials(ctx context.Context, credentials auth.Credentials) error {
	return errors.New("Feature not implemented")
}
