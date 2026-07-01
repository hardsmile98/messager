package server

import (
	"chat/internal/config"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/hardsmile98/messager/sdk/chat/v1"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Dependencies struct {
	Config *config.Config
}

func RunGrpcServer(port string, service pb.ChatServiceServer, deps Dependencies, cleanup func()) error {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	var serverOptions []grpc.ServerOption

	if deps.Config.TLSCertFile != "" && deps.Config.TLSKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(deps.Config.TLSCertFile, deps.Config.TLSKeyFile)

		if err != nil {
			return fmt.Errorf("load tls certificate: %w", err)
		}

		serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		})))
	}

	srv := grpc.NewServer(serverOptions...)

	pb.RegisterChatServiceServer(srv, service)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		slog.Info("shutting down gRPC server")
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		srv.GracefulStop()

		if cleanup != nil {
			cleanup()
		}
	}()

	slog.Info("chat service started", "addr", lis.Addr().String())

	return srv.Serve(lis)
}
