package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/memstorage"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")
	userRepo := memstorage.GetUserRepository()
	lis, err := net.Listen("tcp", "localhost:"+runtimeViper.GetStringSlice("userservice.ports")[0])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", "localhost:"+runtimeViper.GetStringSlice("userservice.ports")[0])
	pb.RegisterUserServiceServer(s, user.GetUserServiceServer(&userRepo))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
