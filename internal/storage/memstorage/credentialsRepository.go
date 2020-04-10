package memstorage

import (
	"errors"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth"
)

type credentialsRepository struct {
	storage *memoryStorage
}

type credentialsEntry struct {
	mu          sync.Mutex // protects credentials
	credentials *auth.Credentials
}

func (cr *credentialsRepository) CreateCredentials(credentials auth.Credentials) error {
	cr.storage.credentialsRWMu.Lock()
	defer cr.storage.credentialsRWMu.Unlock()
	cr.storage.credentials[credentials.Email] = credentials.Password
	return nil
}
func (cr *credentialsRepository) GetCredentials(credentials auth.Credentials) (auth.Credentials, error) {
	cr.storage.credentialsRWMu.Lock()
	defer cr.storage.credentialsRWMu.Unlock()
	pw, exists := cr.storage.credentials[credentials.Email]
	if !exists {
		return auth.Credentials{}, errors.New("username not found")
	}
	return auth.Credentials{credentials.Email, pw}, nil
}
func (cr *credentialsRepository) UpdateCredentials(credentials auth.Credentials) error {
	return errors.New("Feature not implemented")
}
func (cr *credentialsRepository) DeleteCredentials(credentials auth.Credentials) error {
	return errors.New("Feature not implemented")
}
