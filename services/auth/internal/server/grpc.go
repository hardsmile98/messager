package server

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"auth/internal/config"
	"auth/internal/repository"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Dependencies struct {
	Config           *config.Config
	RefreshTokenRepo repository.RefreshTokenRepository
}

func RunGrpcServer(port string, service pb.AuthServiceServer, deps Dependencies, cleanup func()) error {
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

	pb.RegisterAuthServiceServer(srv, service)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		slog.Info("shutting down gRPC server")
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		srv.GracefulStop()

		if cleanup != nil {
			cleanup()
		}
	}()

	slog.Info("auth service started", "addr", lis.Addr().String())

	return srv.Serve(lis)
}
