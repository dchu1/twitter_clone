package storage

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

type AuthRepository struct {
}

func (s *AuthRepository) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {
	result := make(chan *pb.IsAuthenticated, 1)
	errorchan := make(chan error, 1)

	go func() {
		UsersCredRWmu.RLock()
		defer UsersCredRWmu.RUnlock()
		if UsersCred[user.Username] == user.Password {
			result <- &pb.IsAuthenticated{Authenticated: true}
		} else {
			result <- &pb.IsAuthenticated{Authenticated: false}
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

func (s *AuthRepository) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	result := make(chan *pb.Void, 1)
	errorchan := make(chan error, 1)

	go func() {
		UsersCredRWmu.Lock()
		defer UsersCredRWmu.Unlock()
		UsersCred[user.Username] = user.Password
		result <- nil
	}()

	select {
	case res := <-result:
		return res, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		go func() {
			select {
			case <-result:
				delete(UsersCred, user.Username)
				return
			case <-errorchan:
				return
			}
		}()
		return nil, ctx.Err()
	}
}

func (s *AuthRepository) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	result := make(chan *pb.AuthToken, 1)
	errorchan := make(chan error, 1)

	go func() {
		SessionManagerRWmu.Lock()
		defer SessionManagerRWmu.Unlock()
		sessionId := generateSessionId()
		SessionManager[sessionId] = user.UserId
		result <- &pb.AuthToken{Token: sessionId}
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

func (s *AuthRepository) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	result := make(chan *pb.Void, 1)
	errorchan := make(chan error, 1)

	go func() {
		SessionManagerRWmu.Lock()
		defer SessionManagerRWmu.Unlock()
		delete(SessionManager, sess.Token)
		result <- nil
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

func (s *AuthRepository) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {
	result := make(chan *pb.UserId, 1)
	errorchan := make(chan error, 1)

	go func() {
		SessionManagerRWmu.RLock()
		defer SessionManagerRWmu.RUnlock()
		uid, exists := SessionManager[sess.Token]
		if !exists {
			errorchan <- errors.New("invalid token")
			return
		}
		result <- &pb.UserId{UserId: uid}
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
