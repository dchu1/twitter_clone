package memstorage

import (
	"sync"
)

var CredentialsStorage *credentialsStorage

type credentialsStorage struct {
	credentialsRWMu sync.RWMutex // protects the credentials map
	credentials     map[string]*credentialsEntry
}

func NewCredentialsStorage() *credentialsStorage {
	return &credentialsStorage{sync.RWMutex{}, make(map[string]*credentialsEntry)}
}

func InitCredentialsStorage() {
	CredentialsStorage = NewCredentialsStorage()
}
