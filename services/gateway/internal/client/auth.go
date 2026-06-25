package client

import (
	"fmt"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialAuth(url string) (*grpc.ClientConn, pb.AuthServiceClient, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("auth gRPC client: %w", err)
	}

	return conn, pb.NewAuthServiceClient(conn), nil
}
