package handlers

import (
	"log"

	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"google.golang.org/grpc"
)

const (
	authServerAddress = "localhost:50051"
)

var AuthClient authpb.AuthenticationClient

//Register all the grpc clients
func RegisterClients() {
	conn, err := grpc.Dial(authServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()
	AuthClient = authpb.NewAuthenticationClient(conn)
}

func init() {
	RegisterClients()
}
