package auth_test

import (
	"context"
	"strconv"
	"sync"
	"testing"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/server"
)

func TestConcurrentCredential(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)

	authServer := server.GetAuthServer()
	credsArray := make([]*pb.UserCredential, numUsers)
	for i := 0; i < numUsers; i++ {
		credsArray[i] = &pb.UserCredential{Username: strconv.Itoa(i), Password: strconv.Itoa(i)}
	}
	for i := 0; i < numUsers; i++ {
		go func(uid int) {
			defer wg.Done()
			_, err := authServer.AddCredential(context.Background(), credsArray[uid])
			if err != nil {
				t.Error(err.Error())
			}
		}(i)
	}
	wg.Wait()

	wg.Add(numUsers)
	for i := 0; i < numUsers; i++ {
		go func(uid int) {
			defer wg.Done()
			temp, err := authServer.CheckAuthentication(context.Background(), credsArray[uid])
			if err != nil {
				t.Error(err.Error())
			}
			if !temp.Authenticated {
				t.Error("Not Authenticated")
			}
		}(i)
	}
}

func TestConcurrentAuthToken(t *testing.T) {
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	numUsers := 100
	wg.Add(numUsers)
	authServer := server.GetAuthServer()
	authArr := make([]*pb.AuthToken, numUsers)
	for i := 0; i < numUsers; i++ {
		go func(uid uint64) {
			defer wg.Done()
			temp, err := authServer.GetAuthToken(context.Background(), &pb.UserId{UserId: uid})
			if err != nil {
				t.Error(err)
			}
			mu.Lock()
			authArr[uid] = temp
			mu.Unlock()
		}(uint64(i))
	}
	wg.Wait()
	wg.Add(numUsers)
	for i := 0; i < numUsers; i++ {
		go func(uid uint64) {
			defer wg.Done()
			_, err := authServer.GetUserId(context.Background(), authArr[uid])
			if err != nil {
				t.Error(err)
			}
		}(uint64(i))
	}
	wg.Wait()
	wg.Add(numUsers)
	for i := 0; i < numUsers; i++ {
		go func(uid uint64) {
			defer wg.Done()
			_, err := authServer.RemoveAuthToken(context.Background(), authArr[uid])
			if err != nil {
				t.Error(err)
			}
		}(uint64(i))
	}
	wg.Wait()
	wg.Add(numUsers)
	for i := 0; i < numUsers; i++ {
		go func(uid uint64) {
			defer wg.Done()
			_, err := authServer.GetUserId(context.Background(), authArr[uid])
			if err == nil {
				t.Error("user still authenticated")
			}
		}(uint64(i))
	}
}

func GetUserId(t *testing.T) {
	//authServer := server.GetAuthServer()
}
