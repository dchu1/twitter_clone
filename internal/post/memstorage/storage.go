package memstorage

import (
	"errors"
	"sync"
	"time"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/golang/protobuf/ptypes"
)

var PostStorage *postStorage

type postStorage struct {
	postsRWMu sync.RWMutex // protects posts map
	postIDMu  sync.Mutex   // protects postID counter
	posts     map[uint64]*postEntry
	postID    uint64
}

type postEntry struct {
	mu   sync.RWMutex // protects Post
	post *pb.Post
}

func NewPostStorage() *postStorage {
	return &postStorage{sync.RWMutex{}, sync.Mutex{}, make(map[uint64]*postEntry), 1}
}

func (storage *postStorage) createPost(p *pb.Post, result chan uint64, errorchan chan error) {
	storage.postsRWMu.Lock()
	defer storage.postsRWMu.Unlock()
	postEntry := new(postEntry)
	p.PostID = storage.generatePostId()
	postEntry.post = p
	postEntry.post.Timestamp, _ = ptypes.TimestampProto(time.Now())
	storage.posts[p.PostID] = postEntry
	result <- p.PostID
}

func (storage *postStorage) getPost(postID uint64, result chan *pb.Post, errorchan chan error) {
	storage.postsRWMu.RLock()
	defer storage.postsRWMu.RUnlock()
	postEntry, exists := storage.posts[postID]
	if !exists {
		errorchan <- errors.New("user not found")
	} else {
		p := *postEntry.post
		result <- &p
	}
}

func (storage *postStorage) getPosts(postIDs []uint64, result chan []*pb.Post, errorchan chan error) {
	storage.postsRWMu.RLock()
	defer storage.postsRWMu.RUnlock()
	postArr := make([]*pb.Post, 0, len(postIDs))
	for _, v := range postIDs {
		postEntry, _ := storage.posts[v]
		postArr = append(postArr, postEntry.post)
	}
	result <- postArr
}

func (storage *postStorage) getPostsByAuthor(userIDs []uint64, result chan []*pb.Post, errorchan chan error) {
	storage.postsRWMu.RLock()
	defer storage.postsRWMu.RUnlock()
	postArr := make([]*pb.Post, 0, len(userIDs)*100)
	for _, v := range storage.posts {
		v.mu.RLock()
		for _, u := range userIDs {
			if v.post.UserId == u {
				postArr = append(postArr, v.post)
				break
			}
		}
		v.mu.RUnlock()
	}
	result <- postArr
}

func (storage *postStorage) deletePost(postID uint64, errorchan chan error, buffer chan *postEntry) {
	storage.postsRWMu.Lock()
	defer storage.postsRWMu.Unlock()
	postEntry, exists := storage.posts[postID]
	if !exists {
		errorchan <- errors.New("post not exist")
		return
	}
	delete(storage.posts, postID)
	buffer <- postEntry
	errorchan <- nil
}

func init() {
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
