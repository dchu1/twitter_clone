package post_test

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	postmemstorage "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/memstorage"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	usermemstorage "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/memstorage"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

func TestCreatePost(t *testing.T) {
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo := postpb.Post{Message: "testMessage", UserId: user.UserId}
	postID, _ := postApp.CreatePost(context.Background(), &postInfo)
	post, err := postApp.GetPost(context.Background(), postID)
	if err != nil {
		t.Error(err)
	}
	if post.Message != "testMessage" {
		t.Error("Post message not set properly")
	}
	if post.UserId != user.UserId {
		t.Error("Post UserId not set properly")
	}
}
func TestConcurrentCreatePost(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)

	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)

	for post := 0; post < numPosts; post++ {
		go func(post int) {
			defer wg.Done()
			message := "TestMessage " + strconv.Itoa(post)
			postInfo := postpb.Post{Message: message, UserId: user.UserId}
			postApp.CreatePost(context.Background(), &postInfo)
		}(post)
	}
	wg.Wait()
	userArray := []uint64{user.UserId}
	userIDs := postpb.UserIDs{UserIDs: userArray}
	post, err := postApp.GetPostsByAuthors(context.Background(), &userIDs)
	if err != nil {
		t.Error(err)
	}
	if len(post.Posts) != 100 {
		t.Error("Not all posts added")
	}
}
func TestContextCreatePost(t *testing.T) {
	// Create a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo := postpb.Post{Message: "testMessage", UserId: user.UserId}
	postID, err := postApp.CreatePost(ctx, &postInfo)
	if err == nil {
		t.Error("Context error not thrown even after cancelling the context")
	}
	if postID.PostID != 0 {
		t.Error("post created even after cancelling context")
	}
}

func TestGetPost(t *testing.T) {
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	postID1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	postID2, _ := postApp.CreatePost(context.Background(), &postInfo2)
	post1, err1 := postApp.GetPost(context.Background(), postID1)
	post2, err2 := postApp.GetPost(context.Background(), postID2)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil {
		t.Error(err2)
	}
	if post2.Message != "testMessage2" {
		t.Error("Post message not set properly")
	}
	if post2.UserId != user.UserId {
		t.Error("Post UserId not set properly")
	}
	if post1.Message != "testMessage1" {
		t.Error("Post message not set properly")
	}
	if post1.UserId != user.UserId {
		t.Error("Post UserId not set properly")
	}
}
func TestConcurrentGetPost(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Post
	postListmu := sync.Mutex{}

	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postID1, _ := postApp.CreatePost(context.Background(), &postInfo1)

	for post := 0; post < numPosts; post++ {
		go func(post int) {
			defer wg.Done()
			postListmu.Lock()
			defer postListmu.Unlock()
			postValue, _ := postApp.GetPost(context.Background(), postID1)
			postList = append(postList, postValue)
		}(post)
	}
	wg.Wait()

	for _, post := range postList {

		if post.UserId != 1 {
			t.Error("UserID not set properly in the post")
		}
		if post.Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.PostID != 1 {
			t.Error("PostID not set properly in the post")
		}

	}
}
func TestContextGetPost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postID1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	post1, err1 := postApp.GetPost(ctx, postID1)
	if post1 != nil {
		t.Error("post created even after cancelling context")
	}
	if err1 == nil {
		t.Error("Context error not thrown")
	}
}

func TestGetPosts(t *testing.T) {
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	post1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	post2, _ := postApp.CreatePost(context.Background(), &postInfo2)
	postArray := []uint64{post1.PostID, post2.PostID}
	postIDs := postpb.PostIDs{PostIDs: postArray}
	post, err := postApp.GetPosts(context.Background(), &postIDs)
	if err != nil {
		t.Error(err)
	}
	if post.Posts[0].Message != "testMessage2" {
		t.Error("Message not propert for post 2")
	}
	if post.Posts[0].PostID != post2.PostID {
		t.Error("postID not propert for post 2")
	}
	if post.Posts[0].UserId != user.UserId {
		t.Error("userID not propert for post 2")
	}
	if post.Posts[1].Message != "testMessage1" {
		t.Error("Message not propert for post 1")
	}
	if post.Posts[1].PostID != post1.PostID {
		t.Error("postID not propert for post 1")
	}
	if post.Posts[1].UserId != user.UserId {
		t.Error("userID not propert for post 1")
	}
}
func TestConcurrentGetPosts(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Posts
	postListmu := sync.Mutex{}

	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postID1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	postArray := []uint64{postID1.PostID}
	postIDs := postpb.PostIDs{PostIDs: postArray}
	for post := 0; post < numPosts; post++ {
		go func(post int) {
			defer wg.Done()
			postListmu.Lock()
			defer postListmu.Unlock()
			postValue, _ := postApp.GetPosts(context.Background(), &postIDs)
			postList = append(postList, postValue)
		}(post)
	}
	wg.Wait()
	if len(postList) != 100 {
		t.Error("Not all posts received")
	}
	for _, post := range postList {

		if post.Posts[0].UserId != 1 {
			t.Error("UserID not set properly in the post")
		}
		if post.Posts[0].Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.Posts[0].PostID != 1 {
			t.Error("PostID not set properly in the post")
		}

	}
}
func TestContextGetPosts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	post1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	post2, _ := postApp.CreatePost(context.Background(), &postInfo2)
	postArray := []uint64{post1.PostID, post2.PostID}
	postIDs := postpb.PostIDs{PostIDs: postArray}
	post, err := postApp.GetPosts(ctx, &postIDs)
	if post.Posts != nil {
		t.Error("Posts returned even after cancelling the context")
	}
	if err == nil {
		t.Error("Context cancelled still error not thrown")
	}
}

func TestGetPostsByAuthor(t *testing.T) {
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	post1, _ := postApp.CreatePost(context.Background(), &postInfo1)
	post2, _ := postApp.CreatePost(context.Background(), &postInfo2)
	userArray := []uint64{user.UserId}
	userIDs := postpb.UserIDs{UserIDs: userArray}
	post, err := postApp.GetPostsByAuthors(context.Background(), &userIDs)
	if err != nil {
		t.Error(err)
	}
	if post.Posts[0].Message != "testMessage2" {
		t.Error("Message not proper for post 2")
	}
	if post.Posts[0].PostID != post2.PostID {
		t.Error("postID not proper for post 2")
	}
	if post.Posts[0].UserId != user.UserId {
		t.Error("userID not proper for post 2")
	}
	if post.Posts[1].Message != "testMessage1" {
		t.Error("Message not propert for post 1")
	}
	if post.Posts[1].PostID != post1.PostID {
		t.Error("postID not propert for post 1")
	}
	if post.Posts[1].UserId != user.UserId {
		t.Error("userID not propert for post 1")
	}
}
func TestConcurrentGetPostsByAuthor(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Posts
	postListmu := sync.Mutex{}

	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postApp.CreatePost(context.Background(), &postInfo1)
	userArray := []uint64{user.UserId}
	userIDs := postpb.UserIDs{UserIDs: userArray}

	for post := 0; post < numPosts; post++ {
		go func(post int) {
			defer wg.Done()
			postListmu.Lock()
			defer postListmu.Unlock()
			postValue, _ := postApp.GetPostsByAuthors(context.Background(), &userIDs)
			postList = append(postList, postValue)
		}(post)
	}
	wg.Wait()
	if len(postList) != 100 {
		t.Error("Not all posts received")
	}
	for _, post := range postList {

		if post.Posts[0].UserId != 1 {
			t.Error("UserID not set properly in the post")
		}
		if post.Posts[0].Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.Posts[0].PostID != 1 {
			t.Error("PostID not set properly in the post")
		}

	}
}
func TestContextGetPostsByAuthor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postRepo := postmemstorage.GetPostRepository()
	postApp := post.GetPostServiceServer(&postRepo)
	userRepo := usermemstorage.GetUserRepository()
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	postApp.CreatePost(context.Background(), &postInfo1)
	postApp.CreatePost(context.Background(), &postInfo2)
	userArray := []uint64{user.UserId}
	userIDs := postpb.UserIDs{UserIDs: userArray}
	post, err := postApp.GetPostsByAuthors(ctx, &userIDs)
	if post.Posts != nil {
		t.Error("post returned even after cancelling the context")
	}
	if err == nil {
		t.Error("Context cancelled but the error is not thrown")
	}
}
