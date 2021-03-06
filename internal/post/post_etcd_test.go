package post_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	postetcd "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/storage/etcd"
	postmemstorage "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/storage/memstorage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	useretcd "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/storage/etcd"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

func TestCreatePostEtcd(t *testing.T) {
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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
func TestConcurrentCreatePostEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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
func TestContextCreatePostEtcd(t *testing.T) {
	// Create a new context, with its cancellation function
	// from the original context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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
func TestContextTimeoutCreatePostEtcd(t *testing.T) {
	// Create a new context, with its cancellation function
	// from the original context
	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	testPostRepo := postmemstorage.NewTestPostRepository(postRepo)
	postApp := post.GetPostServiceServer(&testPostRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo := postpb.Post{Message: "testMessage", UserId: user.UserId}
	postID1, posterr := postApp.CreatePost(ctx, &postInfo)

	ctx = context.Background()
	_, err := postApp.GetPost(ctx, &postpb.PostID{PostID: uint64(1)})
	if err == nil {
		t.Error(postID1)
		t.Error(posterr)
		t.Error("post still exists")
	}
}

func TestGetPostEtcd(t *testing.T) {
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestConcurrentGetPostEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Post
	postListmu := sync.Mutex{}

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

		if post.UserId != user.UserId {
			t.Error("UserID not set properly in the post")
		}
		if post.Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.PostID != postID1.PostID {
			t.Error("PostID not set properly in the post")
		}

	}
}

func TestContextGetPostEtcd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestContextTimeoutGetPostEtcd(t *testing.T) {
	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	testPostRepo := postmemstorage.NewTestPostRepository(postRepo)
	postApp := post.GetPostServiceServer(&testPostRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestGetPostsEtcd(t *testing.T) {
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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
func TestConcurrentGetPostsEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Posts
	postListmu := sync.Mutex{}

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

		if post.Posts[0].UserId != user.UserId {
			t.Error("UserID not set properly in the post")
		}
		if post.Posts[0].Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.Posts[0].PostID != postID1.PostID {
			t.Error("PostID not set properly in the post")
		}

	}
}
func TestContextGetPostsEtcd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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
func TestContextTimeoutGetPostsEtcd(t *testing.T) {
	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	testPostRepo := postmemstorage.NewTestPostRepository(postRepo)
	postApp := post.GetPostServiceServer(&testPostRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestGetPostsByAuthorEtcd(t *testing.T) {
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestConcurrentGetPostsByAuthorEtcd(t *testing.T) {
	var wg sync.WaitGroup
	numPosts := 100
	wg.Add(numPosts)
	var postList []*postpb.Posts
	postListmu := sync.Mutex{}

	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	post1, _ := postApp.CreatePost(context.Background(), &postInfo1)
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

		if post.Posts[0].UserId != user.UserId {
			t.Error("UserID not set properly in the post")
		}
		if post.Posts[0].Message != "testMessage1" {
			t.Error("Message not set properly in the post")
		}
		if post.Posts[0].PostID != post1.PostID {
			t.Error("PostID not set properly in the post")
		}

	}
}

func TestContextGetPostsByAuthorEtcd(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	postApp := post.GetPostServiceServer(&postRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
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

func TestContextTimeoutGetPostsByAuthorEtcd(t *testing.T) {
	duration := 15 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), duration)
	postStorage, _ := postetcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer postStorage.Close()
	postRepo := postetcd.NewPostRepository(postStorage)
	testPostRepo := postmemstorage.NewTestPostRepository(postRepo)
	postApp := post.GetPostServiceServer(&testPostRepo)

	userStorage, _ := useretcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	defer userStorage.Close()
	userRepo := useretcd.NewUserRepository(userStorage)
	userApp := user.GetUserServiceServer(&userRepo)

	userInfo := userpb.AccountInformation{FirstName: "test1", LastName: "test2", Email: "test@nyu.edu"}
	user, _ := userApp.CreateUser(context.Background(), &userInfo)
	postInfo1 := postpb.Post{Message: "testMessage1", UserId: user.UserId}
	postInfo2 := postpb.Post{Message: "testMessage2", UserId: user.UserId}
	go postApp.CreatePost(context.Background(), &postInfo1)
	go postApp.CreatePost(context.Background(), &postInfo2)
	userArray := []uint64{user.UserId}
	userIDs := postpb.UserIDs{UserIDs: userArray}
	time.Sleep(1 * time.Second)
	post, err := postApp.GetPostsByAuthors(ctx, &userIDs)
	if post.Posts != nil {
		t.Error("post returned even after cancelling the context")
	}
	if err == nil {
		t.Error("Context cancelled but the error is not thrown")
	}
}
