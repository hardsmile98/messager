package client

import (
	"fmt"

	pb "github.com/hardsmile98/messager/sdk/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialChat(url string) (*grpc.ClientConn, pb.ChatServiceClient, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("chat gRPC client: %w", err)
	}

	return conn, pb.NewChatServiceClient(conn), nil
}
