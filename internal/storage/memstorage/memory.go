package memstorage

import (
	"errors"
	"sync"
)

type memoryStorage struct {
	usersRWMu       sync.RWMutex // protects users map
	postsRWMu       sync.RWMutex // protects posts map
	userIDMu        sync.Mutex   // protects userID counter
	postIDMu        sync.Mutex   // protects postID counter
	credentialsRWMu sync.RWMutex // protects the credentials map
	credentials     map[string]*credentialsEntry
	users           map[uint64]*userEntry
	posts           map[uint64]*postEntry
	userID          uint64
	postID          uint64
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{sync.RWMutex{}, sync.RWMutex{}, sync.Mutex{}, sync.Mutex{}, sync.RWMutex{}, make(map[string]*credentialsEntry), make(map[uint64]*userEntry), make(map[uint64]*postEntry), 0, 0}
}

// generateUserId gets a value from the userId counter, then increments the counter
func (storage *memoryStorage) generateUserId() uint64 {
	storage.userIDMu.Lock()
	defer storage.userIDMu.Unlock()
	uid := storage.userID
	storage.userID++
	return uid
}

// generatePostId gets a value from the postId counter, then increments the counter
func (storage *memoryStorage) generatePostId() uint64 {
	storage.postIDMu.Lock()
	defer storage.postIDMu.Unlock()
	uid := storage.postID
	storage.postID++
	return uid
}

// getUserEntry is a function for getting a UserEntry object (not a clone)
func (storage *memoryStorage) getUserEntry(userID uint64) (*userEntry, error) {
	storage.usersRWMu.RLock()
	defer storage.usersRWMu.RUnlock()
	u, exists := storage.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}
