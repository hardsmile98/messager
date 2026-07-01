package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantRepository interface {
	IsParticipant(ctx context.Context, chatID, userID string) (bool, error)
}

type ParticipantRepo struct {
	pool *pgxpool.Pool
}

func NewParticipantRepo(pool *pgxpool.Pool) *ParticipantRepo {
	return &ParticipantRepo{pool}
}

func (r *ParticipantRepo) IsParticipant(ctx context.Context, chatID, userID string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2)
	`

	err := r.pool.QueryRow(ctx, query, chatID, userID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
