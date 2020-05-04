package auth_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	server "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/service"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/storage/etcd"
)

func TestEtcdConcurrentCredential(t *testing.T) {
	var wg sync.WaitGroup
	numUsers := 100
	wg.Add(numUsers)

	client, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	authServer := server.GetEtcdAuthServer(client)
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

func TestEtcdContextTimeoutCredential(t *testing.T) {
	var err error
	client, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	authServer := server.GetTestEtcdAuthServer(client)
	errchan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	cred := &pb.UserCredential{Username: "test", Password: "test"}
	go func() {
		_, err := authServer.AddCredential(ctx, cred)
		errchan <- err
	}()
	time.Sleep(1 * time.Second)
	cancel()
	err = <-errchan
	if err == nil {
		t.Error("No error returned")
	}
	result, err := authServer.CheckAuthentication(context.Background(), cred)
	if err != nil {
		t.Error(err.Error())
	}
	if result.Authenticated {
		t.Error("User authenticated")
	}
}

func TestEtcdContextTimeoutAuthToken(t *testing.T) {
	var err error
	client, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	authServer := server.GetTestEtcdAuthServer(client)
	resultchan := make(chan *pb.AuthToken)
	errchan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		temp, err := authServer.GetAuthToken(ctx, &pb.UserId{UserId: uint64(1)})
		if err != nil {
			errchan <- err
		}
		resultchan <- temp

	}()
	time.Sleep(1 * time.Second)
	cancel()
	select {
	case <-resultchan:
	case err = <-errchan:
	}
	if err == nil {
		t.Error("No error returned")
	}

}

func TestEtcdContextTimeoutRemoveAuthToken(t *testing.T) {
	var err error
	client, _ := etcd.NewClient([]string{"http://localhost:2379", "http://localhost:22379", "http://localhost:32379"})
	authServer := server.GetTestEtcdAuthServer(client)
	token, err := authServer.GetAuthToken(context.Background(), &pb.UserId{UserId: uint64(1)})

	errchan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, err := authServer.RemoveAuthToken(ctx, token)
		errchan <- err
	}()
	time.Sleep(1 * time.Second)
	cancel()
	err = <-errchan
	if err == nil {
		t.Error("No error returned")
	}

	_, err = authServer.GetUserId(context.Background(), token)
	if err != nil {
		t.Error("Token removed")
	}
}
