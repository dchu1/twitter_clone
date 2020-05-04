package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/storage/etcd"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/storage/memstorage"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	// Read Config
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")
	var userRepo user.UserRepository

	// Initialize storage
	switch runtimeViper.GetString("userservice.storage") {
	case "etcd":
		etcdCli, err := etcd.NewClient(runtimeViper.GetStringSlice("userservice.etcdcluster"))
		if err != nil {
			log.Fatal(err)
		}
		userRepo = etcd.NewUserRepository(etcdCli)
	case "memory":
		userRepo = memstorage.GetUserRepository()
	default:
		log.Fatal("Could not read config for storage")
	}

	// Start server
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
