package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"google.golang.org/grpc"

	db "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage"
)

const (
	port = ":50051"
)

type authServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *authServer) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {

	if db.UsersCred[user.Username] == user.Password {
		return &pb.IsAuthenticated{Authenticated: true}, nil
	}
	return &pb.IsAuthenticated{Authenticated: false}, nil
}

func (s *authServer) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	db.UsersCred[user.Username] = user.Password
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthenticationServer(s, &authServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
