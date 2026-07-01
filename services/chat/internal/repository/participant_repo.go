package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ParticipantRepository interface{}

type ParticipantRepo struct {
	pool *pgxpool.Pool
}

func NewParticipantRepo(pool *pgxpool.Pool) *ParticipantRepo {
	return &ParticipantRepo{pool}
}
