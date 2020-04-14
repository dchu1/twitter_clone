package memstorage

import (
	"errors"
	"sync"
)

var UserStorage *userStorage

type userStorage struct {
	usersRWMu sync.RWMutex // protects users map
	userIDMu  sync.Mutex   // protects userID counter
	users     map[uint64]*userEntry
	userID    uint64
}

func NewUserStorage() *userStorage {
	return &userStorage{sync.RWMutex{}, sync.Mutex{}, make(map[uint64]*userEntry), 1}
}

func init() {
	UserStorage = NewUserStorage()
}

// generateUserId gets a value from the userId counter, then increments the counter
func (storage *userStorage) generateUserId() uint64 {
	storage.userIDMu.Lock()
	defer storage.userIDMu.Unlock()
	uid := storage.userID
	storage.userID++
	return uid
}

// getUserEntry is a function for getting a UserEntry object (not a clone)
func (storage *userStorage) getUserEntry(userID uint64) (*userEntry, error) {
	storage.usersRWMu.RLock()
	defer storage.usersRWMu.RUnlock()
	u, exists := storage.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}
