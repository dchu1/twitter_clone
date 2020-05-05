package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/feed"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/feed/feedpb"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	// Read Config
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")

	// Get our clients
	postServerAddress := "localhost:" + runtimeViper.GetStringSlice("postservice.ports")[0]
	userServerAddress := "localhost:" + runtimeViper.GetStringSlice("userservice.ports")[0]

	conn, err := grpc.Dial(userServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	userServiceClient := userpb.NewUserServiceClient(conn)

	conn, err = grpc.Dial(postServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	postServiceClient := postpb.NewPostServiceClient(conn)

	// Start server
	lis, err := net.Listen("tcp", "localhost:"+runtimeViper.GetStringSlice("feedservice.ports")[0])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", "localhost:"+runtimeViper.GetStringSlice("feedservice.ports")[0])
	feedpb.RegisterFeedServiceServer(s, feed.NewFeedServiceServer(postServiceClient, userServiceClient))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
