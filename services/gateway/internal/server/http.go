package server

import (
	"context"
	"fmt"
	"gateway/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const shutdownTimeout = 10 * time.Second

func RunHTTPServer(conf *config.Config) error {
	authConnection, err := grpc.NewClient(
		conf.AuthGRPCURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("auth gRPC client: %w", err)
	}

	defer func() {
		if err := authConnection.Close(); err != nil {
			slog.Error("failed to close auth gRPC connection", "error", err)
		}
	}()

	authClient := pb.NewAuthServiceClient(authConnection)

	srv := &http.Server{
		Addr:              ":" + conf.Port,
		Handler:           newRouter(authClient, conf),
		ReadTimeout:       conf.HTTPReadTimeout,
		WriteTimeout:      conf.HTTPWriteTimeout,
		IdleTimeout:       conf.HTTPIdleTimeout,
		ReadHeaderTimeout: conf.HTTPReadHeaderTimeout,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		slog.Info("shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("http shutdown", "error", err)
		}
	}()

	slog.Info("gateway started", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}
