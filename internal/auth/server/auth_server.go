package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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

func (s *authServer) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	sessionId := generateSessionId()
	db.SessionManager[sessionId] = user.UserId
	return &pb.AuthToken{Token: sessionId}, nil
}

func (s *authServer) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	delete(db.SessionManager, sess.Token)
	return nil, nil
}

func (s *authServer) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {

	return &pb.UserId{UserId: db.SessionManager[sess.Token]}, nil
}

func generateSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fmt.Println("Server running on port", port)
	pb.RegisterAuthenticationServer(s, &authServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
