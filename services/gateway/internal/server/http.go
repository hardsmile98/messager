package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway/internal/client"
	"gateway/internal/config"
	httptransport "gateway/internal/transport/http"
)

const shutdownTimeout = 10 * time.Second

func RunHTTPServer(conf *config.Config) error {
	authConn, authClient, err := client.DialAuth(conf.AuthGRPCURL)

	if err != nil {
		return err
	}

	defer func() {
		if err := authConn.Close(); err != nil {
			slog.Error("failed to close auth gRPC connection", "error", err)
		}
	}()

	srv := &http.Server{
		Addr:              ":" + conf.Port,
		Handler:           httptransport.NewRouter(authClient, conf),
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
