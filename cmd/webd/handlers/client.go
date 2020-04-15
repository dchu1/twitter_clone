package handlers

import (
	"fmt"
	"log"

	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	postpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var AuthClient authpb.AuthenticationClient
var UserServiceClient userpb.UserServiceClient
var PostServiceClient postpb.PostServiceClient

//Register all the grpc clients
func RegisterClients() {
	fmt.Printf("%+v", viper.AllSettings())
	authServerAddress := "localhost:" + viper.GetStringSlice("authservice.ports")[0]
	postServerAddress := "localhost:" + viper.GetStringSlice("postservice.ports")[0]
	userServerAddress := "localhost:" + viper.GetStringSlice("userservice.ports")[0]

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
}

// func init() {
// 	RegisterClients()
// }
