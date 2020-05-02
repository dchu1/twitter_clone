package memstorage

import (
	"context"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

type testAuthRepository struct {
	authRepo auth.AuthRepository
}

func GetTestAuthRepository(a auth.AuthRepository) auth.AuthRepository {
	return &testAuthRepository{a}
}

func (s *testAuthRepository) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {
	time.Sleep(time.Second * 5)
	return s.authRepo.CheckAuthentication(ctx, user)

}

func (s *testAuthRepository) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	time.Sleep(time.Second * 5)
	return s.authRepo.AddCredential(ctx, user)
}

func (s *testAuthRepository) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	time.Sleep(time.Second * 5)
	return s.authRepo.GetAuthToken(ctx, user)
}

func (s *testAuthRepository) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	time.Sleep(time.Second * 5)
	return s.authRepo.RemoveAuthToken(ctx, sess)
}

func (s *testAuthRepository) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {
	time.Sleep(time.Second * 5)
	return s.authRepo.GetUserId(ctx, sess)
}
