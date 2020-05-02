package etcd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"go.etcd.io/etcd/clientv3"
)

type authRepository struct {
	storage *clientv3.Client
}

func GetAuthRepository(client *clientv3.Client) auth.AuthRepository {
	return &authRepository{client}
	// return new(authRepository)
}

func (auth *authRepository) CheckAuthentication(ctx context.Context, user *pb.UserCredential) (*pb.IsAuthenticated, error) {
	result := make(chan *pb.IsAuthenticated, 1)
	errorchan := make(chan error, 1)

	go func() {
		// UsersCredRWmu.RLock()
		// defer UsersCredRWmu.RUnlock()
		resp, err := auth.storage.Get(ctx, "UserCred")
		if err != nil {
			errorchan <- err
		} else {
			var usr []byte
			for _, ev := range resp.Kvs {
				usr = ev.Value
			}
			users := make(map[string]interface{})
			json.Unmarshal([]byte(usr), &users)

			if users[user.Username] == user.Password {
				result <- &pb.IsAuthenticated{Authenticated: true}
			} else {
				result <- &pb.IsAuthenticated{Authenticated: false}
			}
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

func (auth *authRepository) AddCredential(ctx context.Context, user *pb.UserCredential) (*pb.Void, error) {
	result := make(chan *pb.Void, 1)
	errorchan := make(chan error, 1)

	go func() {
		// UsersCredRWmu.Lock()
		// defer UsersCredRWmu.Unlock()

		resp, err := auth.storage.Get(ctx, "UserCred")
		if err != nil {
			errorchan <- err
		} else {
			var usr []byte
			for _, ev := range resp.Kvs {
				usr = ev.Value
			}
			users := make(map[string]interface{})
			json.Unmarshal([]byte(usr), &users)
			users[user.Username] = user.Password // add user
			usersJson, err := json.Marshal(users)

			usersStr := string(usersJson)
			_, err = auth.storage.Put(context.Background(), "UserCred", usersStr)
			if err != nil {
				errorchan <- err
			} else {
				result <- nil
			}
		}

	}()

	select {
	case res := <-result:
		return res, nil
	case err := <-errorchan:
		return &pb.Void{}, err
	case <-ctx.Done():
		go func() {
			select {
			case <-result:
				// UsersCredRWmu.Lock()
				// defer UsersCredRWmu.Unlock()

				resp, err := auth.storage.Get(ctx, "UserCred")
				if err != nil {
					errorchan <- err
				} else {
					var usr []byte
					for _, ev := range resp.Kvs {
						usr = ev.Value
					}
					users := make(map[string]interface{})
					json.Unmarshal([]byte(usr), &users)
					delete(users, user.Username) // delete user
					usersJson, err := json.Marshal(users)
					usersStr := string(usersJson)
					_, err = auth.storage.Put(context.Background(), "UserCred", usersStr)
					if err != nil {
						return
					} else {
						return
					}
				}
				return
			case <-errorchan:
				return
			}
		}()
		return &pb.Void{}, ctx.Err()
	}
}

func (auth *authRepository) GetAuthToken(ctx context.Context, user *pb.UserId) (*pb.AuthToken, error) {
	result := make(chan *pb.AuthToken, 1)
	errorchan := make(chan error, 1)

	go func() {
		// SessionManagerRWmu.Lock()
		// defer SessionManagerRWmu.Unlock()
		sessionId := generateSessionId()
		resp, err := auth.storage.Get(ctx, "SessionStore")
		if err != nil {
			errorchan <- err
		} else {
			var sess []byte
			for _, ev := range resp.Kvs {
				sess = ev.Value
			}
			sessions := make(map[string]interface{})
			json.Unmarshal([]byte(sess), &sessions)
			sessions[sessionId] = user.UserId // add session

			sessionsJson, err := json.Marshal(sessions)

			sessionsStr := string(sessionsJson)
			_, err = auth.storage.Put(context.Background(), "SessionStore", sessionsStr)
			if err != nil {
				errorchan <- err
			} else {
				result <- &pb.AuthToken{Token: sessionId}
			}
		}

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

func (auth *authRepository) RemoveAuthToken(ctx context.Context, sess *pb.AuthToken) (*pb.Void, error) {
	errorchan := make(chan error, 1)
	bufferchan := make(chan uint64, 1)

	go func() {
		// SessionManagerRWmu.Lock()
		// defer SessionManagerRWmu.Unlock()

		resp, err := auth.storage.Get(ctx, "SessionStore")
		if err != nil {
			errorchan <- err
		} else {
			var sessMap []byte
			for _, ev := range resp.Kvs {
				sessMap = ev.Value
			}
			sessions := make(map[string]interface{})
			json.Unmarshal([]byte(sessMap), &sessions)

			token, exists := sessions[sess.Token]
			if !exists {
				errorchan <- errors.New("token does not exist")
				return
			}

			delete(sessions, sess.Token) // delete the session
			sessionsJson, err := json.Marshal(sessions)

			sessionsStr := string(sessionsJson)
			_, err = auth.storage.Put(context.Background(), "SessionStore", sessionsStr)
			if err != nil {
				errorchan <- err
			} else {
				bufferchan <- uint64(token.(float64))
			}
		}
		// bufferchan <- token
	}()

	select {
	case err := <-errorchan:
		return &pb.Void{}, err
	case <-ctx.Done():
		go func() {
			select {
			case token := <-bufferchan:
				// SessionManagerRWmu.Lock()
				// defer SessionManagerRWmu.Unlock()
				resp, err := auth.storage.Get(ctx, "SessionStore")
				if err != nil {
					errorchan <- err
				} else {
					var sessMap []byte
					for _, ev := range resp.Kvs {
						sessMap = ev.Value
					}
					sessions := make(map[string]interface{})
					json.Unmarshal([]byte(sessMap), &sessions)
					sessions[sess.Token] = token // add session back
					sessionsJson, _ := json.Marshal(sessions)

					sessionsStr := string(sessionsJson)
					_, _ = auth.storage.Put(context.Background(), "SessionStore", sessionsStr)

				}
				return
			case <-errorchan:
				return
			}
		}()
		return &pb.Void{}, ctx.Err()
	}
}

func (auth *authRepository) GetUserId(ctx context.Context, sess *pb.AuthToken) (*pb.UserId, error) {
	result := make(chan *pb.UserId, 1)
	errorchan := make(chan error, 1)

	go func() {
		// SessionManagerRWmu.RLock()
		// defer SessionManagerRWmu.RUnlock()

		resp, err := auth.storage.Get(ctx, "SessionStore")
		if err != nil {
			errorchan <- err
		} else {
			var sessMap []byte
			for _, ev := range resp.Kvs {
				sessMap = ev.Value
			}
			sessions := make(map[string]interface{})

			json.Unmarshal([]byte(sessMap), &sessions)
			userId, exists := sessions[sess.Token]
			if !exists {
				errorchan <- errors.New("invalid token")
				return
			} else {
				result <- &pb.UserId{UserId: uint64(userId.(float64))}
			}
		}

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
