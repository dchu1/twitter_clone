package memstorage

import (
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
)

var PostStorage *postStorage

type postStorage struct {
	postsRWMu sync.RWMutex // protects posts map
	postIDMu  sync.Mutex   // protects postID counter
	posts     map[uint64]*postEntry
	postID    uint64
}

type postEntry struct {
	mu   sync.Mutex // protects Post
	post *post.Post
}

func NewPostStorage() *postStorage {
	return &postStorage{sync.RWMutex{}, sync.Mutex{}, make(map[uint64]*postEntry), 0}
}

func InitPostStorage() {
	PostStorage = NewPostStorage()
}

// generatePostId gets a value from the postId counter, then increments the counter
func (storage *postStorage) generatePostId() uint64 {
	storage.postIDMu.Lock()
	defer storage.postIDMu.Unlock()
	uid := storage.postID
	storage.postID++
	return uid
}
