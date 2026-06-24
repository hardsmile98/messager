package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	pb "sdk/auth/v1"
	"syscall"

	"google.golang.org/grpc"
)

func RunGrpcServer(port string, service pb.AuthServiceServer) error {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	fmt.Println("Auth service running on port: ", lis.Addr().String())

	srv := grpc.NewServer()

	pb.RegisterAuthServiceServer(srv, service)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		srv.GracefulStop()
	}()

	return srv.Serve(lis)
}
