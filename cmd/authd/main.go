package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/server"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	runtimeViper := viper.New()
	config.NewRuntimeConfig(runtimeViper, ".")
	lis, err := net.Listen("tcp", "localhost:"+runtimeViper.GetStringSlice("authservice.ports")[0])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", "localhost:"+runtimeViper.GetStringSlice("authservice.ports")[0])
	pb.RegisterAuthenticationServer(s, server.GetAuthServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
