package main

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.LoadConfig()

	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBUri)

	if err != nil {
		slog.Error("failed to initialize database pool", "error", err)
		os.Exit(1)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		slog.Error("failed to initialize redis", "error", err)
		pool.Close()
		os.Exit(1)
	}

	userRepo := repository.NewUserRepo(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepo(pool, redisClient)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, cfg)

	cleanup := func() {
		pool.Close()

		if err := redisClient.Close(); err != nil {
			slog.Error("failed to close redis client", "error", err)
		}
	}

	err = server.RunGrpcServer(cfg.GRPCPort, authService, server.Dependencies{
		Config:           cfg,
		RefreshTokenRepo: refreshTokenRepo,
	}, cleanup)

	if err != nil {
		cleanup()
		slog.Error("failed to run gRPC server", "error", err)
	}
}
