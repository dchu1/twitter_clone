package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"

	db "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage"
)

const (
	port = ":50051"
)

type authServer struct {
	pb.UnimplementedAuthenticationServer
}

func (s *authServer) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {
	result := make(chan *pb.IsAuthenticated)
	errorchan := make(chan error)

	go func() {

		if db.UsersCred[user.Username] == user.Password {

			result <- &pb.IsAuthenticated{Authenticated: true}
			errorchan <- nil
		} else {
			result <- &pb.IsAuthenticated{Authenticated: false}
			errorchan <- nil
		}

	}()

	select {
	case auth := <-result:
		return auth, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *authServer) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	result := make(chan *pb.Void)
	errorchan := make(chan error)

	go func() {
		db.UsersCred[user.Username] = user.Password
		result <- nil
		errorchan <- nil
	}()

	select {
	case res := <-result:
		return res, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		delete(db.UsersCred, user.Username)
		return nil, ctx.Err()
	}
}

func (s *authServer) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	result := make(chan *pb.AuthToken)
	errorchan := make(chan error)

	go func() {

		sessionId := generateSessionId()
		db.SessionManager[sessionId] = user.UserId
		result <- &pb.AuthToken{Token: sessionId}
		errorchan <- nil
	}()

	select {
	case token := <-result:
		return token, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *authServer) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	result := make(chan *pb.Void)
	errorchan := make(chan error)

	go func() {
		delete(db.SessionManager, sess.Token)
		result <- nil
		errorchan <- nil
	}()

	select {
	case token := <-result:
		return token, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *authServer) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {
	result := make(chan *pb.UserId)
	errorchan := make(chan error)

	go func() {

		result <- &pb.UserId{UserId: db.SessionManager[sess.Token]}
		errorchan <- nil
	}()

	select {
	case userID := <-result:
		return userID, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func generateSessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func GetAuthServer() *authServer {
	return &authServer{}
}

// func main() {

// 	lis, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	fmt.Println("Server running on port", port)
// 	pb.RegisterAuthenticationServer(s, &authServer{})
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
