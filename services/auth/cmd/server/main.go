package main

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initPool(dbUri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbUri)

	if err != nil {
		return nil, err
	}

	defer pool.Close()

	return pool, nil
}

func initRedis(redisAddr string) (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := redis.Ping(context.Background()).Result()

	return redis, err
}

func main() {
	config, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := initPool(config.DBUri)

	if err != nil {
		log.Fatalf("Failed to initialize database pool: %v", err)
	}

	redis, err := initRedis(config.RedisAddr)

	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	userRepo := repository.NewUserRepo(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepo(pool, redis)

	authService := service.NewAuthService(userRepo, refreshTokenRepo, config)

	err = server.RunGrpcServer(config.GRPCPort, authService)

	if err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
