package memstorage

import (
	"errors"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

// UserStorage is package level userStorage, which can be used instead of instantiating one
var UserStorage *userStorage

type userStorage struct {
	usersRWMu sync.RWMutex // protects users map
	userIDMu  sync.Mutex   // protects userID counter
	users     map[uint64]*userEntry
	userID    uint64
}

type userEntry struct {
	followingRWMu sync.RWMutex // protects following map
	followersRWMu sync.RWMutex // protects followers map
	user          *userpb.User
}

// NewUserStorage returns a new instance of a userStorage
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
