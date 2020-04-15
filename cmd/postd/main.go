package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/memstorage"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")
	postRepo := memstorage.GetPostRepository()
	lis, err := net.Listen("tcp", "localhost:"+runtimeViper.GetStringSlice("postservice.ports")[0])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", "localhost:"+runtimeViper.GetStringSlice("postservice.ports")[0])
	pb.RegisterPostServiceServer(s, post.GetPostServiceServer(&postRepo))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
