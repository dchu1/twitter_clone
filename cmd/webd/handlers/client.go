package handlers

import (
	"log"

	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	feedpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/feed/feedpb"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	signuppb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/signup/signuppb"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var AuthClient authpb.AuthenticationClient
var UserServiceClient userpb.UserServiceClient
var PostServiceClient postpb.PostServiceClient
var FeedServiceClient feedpb.FeedServiceClient
var SignupServiceClient signuppb.SignupServiceClient

//Register all the grpc clients
func RegisterClients() {
	authServerAddress := "localhost:" + viper.GetStringSlice("authservice.ports")[0]
	postServerAddress := "localhost:" + viper.GetStringSlice("postservice.ports")[0]
	userServerAddress := "localhost:" + viper.GetStringSlice("userservice.ports")[0]
	feedServerAddress := "localhost:" + viper.GetStringSlice("feedservice.ports")[0]
	signupServerAddress := "localhost:" + viper.GetStringSlice("signupservice.ports")[0]

	conn, err := grpc.Dial(authServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()
	AuthClient = authpb.NewAuthenticationClient(conn)

	conn, err = grpc.Dial(userServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	UserServiceClient = userpb.NewUserServiceClient(conn)

	conn, err = grpc.Dial(postServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	PostServiceClient = postpb.NewPostServiceClient(conn)

	conn, err = grpc.Dial(feedServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	FeedServiceClient = feedpb.NewFeedServiceClient(conn)

	conn, err = grpc.Dial(signupServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	SignupServiceClient = signuppb.NewSignupServiceClient(conn)
}

// func init() {
// 	RegisterClients()
// }
