package post_test

import (
	"context"
	"log"
	"reflect"
	"testing"

	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"google.golang.org/grpc"
)

const postServerAddress = "localhost:50052"

var PostServiceClient postpb.PostServiceClient

func init() {
	conn, err := grpc.Dial(postServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	PostServiceClient = postpb.NewPostServiceClient(conn)
}
func TestCreatePost(t *testing.T) {
	_, err := PostServiceClient.CreatePost(context.Background(), &postpb.Post{UserId: 0, Message: "Test Message"})
	if err != nil {
		t.Error(err.Error())
	}
	retPost, err := PostServiceClient.GetPost(context.Background(), &postpb.PostID{PostID: 0})
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(retPost.Message, "Test Message") || !reflect.DeepEqual(retPost.PostID, 0) || !reflect.DeepEqual(retPost.UserId, 0) {
		t.Error("Test Failed Posts struct not in sync")
	}
}

func UnaryClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Create a new context with the token and make the first request
	// sleepTime, _ := time.Parse(opts.SleepTime)
	// time.Sleep(sleepTime)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
