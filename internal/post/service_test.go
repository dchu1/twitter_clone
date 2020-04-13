package post_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/memstorage"
)

func TestCreatePost(t *testing.T) {
	storage := memstorage.NewPostStorage()
	app := post.NewService(memstorage.NewPostRepository(storage))
	p := post.Post{0, time.Now(), "Test Message", 0}
	app.CreatePost(nil, p)

	retPost, err := app.GetPost(nil, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(retPost.Message, p.Message) || !reflect.DeepEqual(retPost.PostID, p.PostID) || !reflect.DeepEqual(retPost.UserID, p.UserID) {
		t.Error("Test Failed Posts struct not in sync")
	}
}
