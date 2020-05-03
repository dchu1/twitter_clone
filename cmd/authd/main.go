package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/service"
	etcd "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage/etcd"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {


	// Read Config
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")


	storage := runtimeViper.GetStringSlice("storage.storage")[0]


	// Start server
	lis, err := net.Listen("tcp", "localhost:"+runtimeViper.GetStringSlice("authservice.ports")[0])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", "localhost:"+runtimeViper.GetStringSlice("authservice.ports")[0])
	if storage == "etcd" {
		client, _ := etcd.NewClient([]string{"http://localhost:2379"})
		pb.RegisterAuthenticationServer(s, service.GetEtcdAuthServer(client))
	} else {
		pb.RegisterAuthenticationServer(s, service.GetAuthServer())
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
