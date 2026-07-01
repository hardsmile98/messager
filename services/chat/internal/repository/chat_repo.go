package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ChatRepository interface{}

type ChatRepo struct {
	pool *pgxpool.Pool
}

func NewChatRepo(pool *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{pool}
}
