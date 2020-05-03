package service

import (
	"context"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	etcdstorage "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage/etcd"
	memstorage "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage/memstorage"
	"go.etcd.io/etcd/clientv3"
)

const (
	port = ":50051"
)

type authServer struct {
	authRepository auth.AuthRepository
	pb.UnimplementedAuthenticationServer
}

func (s *authServer) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {
	return s.authRepository.CheckAuthentication(ctx, user)
}

func (s *authServer) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	return s.authRepository.AddCredential(ctx, user)
}

func (s *authServer) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	return s.authRepository.GetAuthToken(ctx, user)
}

func (s *authServer) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	return s.authRepository.RemoveAuthToken(ctx, sess)
}

func (s *authServer) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {
	return s.authRepository.GetUserId(ctx, sess)
}

func GetEtcdAuthServer(client *clientv3.Client) *authServer {
	return &authServer{authRepository: etcdstorage.GetAuthRepository(client)}
}

func GetAuthServer() *authServer {
	return &authServer{authRepository: memstorage.GetAuthRepository()}
}

func GetTestAuthServer() *authServer {
	a := memstorage.GetAuthRepository()
	b := memstorage.GetTestAuthRepository(a)
	return &authServer{authRepository: b}
}

func GetTestEtcdAuthServer(client *clientv3.Client) *authServer {
	a := etcdstorage.GetAuthRepository(client)
	b := etcdstorage.GetTestAuthRepository(a)
	return &authServer{authRepository: b}
}
