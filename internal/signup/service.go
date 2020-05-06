package signup

import (
	"context"

	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	signuppb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/signup/signuppb"
	userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"github.com/golang/protobuf/ptypes/empty"
)

type signupServiceServer struct {
	AuthClient        authpb.AuthenticationClient
	UserServiceClient userpb.UserServiceClient
	signuppb.UnimplementedSignupServiceServer
}

func (s *signupServiceServer) Signup(ctx context.Context, req *signuppb.AccountInformation) (*empty.Empty, error) {
	// Convert the signuppb.AccountInformation to userpb.AccountInformation
	// Ideally I would use the same object, but I don't know how to import protoc files from
	// other protoc files...
	info := &userpb.AccountInformation{FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}

	// There must be some way to generalize this context waiting pattern...
	errorchan := make(chan error, 1)
	go func() {
		_, err := s.UserServiceClient.CreateUser(ctx, info)
		errorchan <- err
	}()

	select {
	case err := <-errorchan:
		if err == nil {
			errorchan2 := make(chan error, 1)
			go func() {
				_, err := s.AuthClient.AddCredential(ctx, &authpb.UserCredential{Username: req.Email, Password: req.Password})
				errorchan2 <- err
			}()
			select {
			case err := <-errorchan2:
				if err != nil {
					go s.UserServiceClient.DeleteUser(context.Background(), info)
				}
				return new(empty.Empty), nil
			case <-ctx.Done():
				go s.UserServiceClient.DeleteUser(context.Background(), info)

				// If ctx.Done(), i'm going to assume that AuthService cleaned up after itself.
				// There is a scenario where ctx.Done() gets called after AuthService has created
				// the entry, but before it has returned. Since we don't have a delete op for
				// AuthService yet, this is unhandled
				go func() {
					err := <-errorchan2
					if err != nil && err != ctx.Err() {
						// TODO Implement delete operation for AuthService
					}
					return
				}()
				return new(empty.Empty), ctx.Err()
			}
		} else {
			return new(empty.Empty), err
		}
	case <-ctx.Done():
		go func() {
			err := <-errorchan
			if err != nil && err != ctx.Err() {
				s.UserServiceClient.DeleteUser(context.Background(), info)
			}
			return
		}()

		return new(empty.Empty), ctx.Err()
	}

}

// NewSignupServiceServer returns a grpc Server for the signup service using the provided clients
func NewSignupServiceServer(asc authpb.AuthenticationClient, usc userpb.UserServiceClient) *signupServiceServer {
	return &signupServiceServer{AuthClient: asc, UserServiceClient: usc}
}
