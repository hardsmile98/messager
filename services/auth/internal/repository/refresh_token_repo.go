package repository

import (
	"auth/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepo struct {
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewRefreshTokenRepo(pool *pgxpool.Pool, redis *redis.Client) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool, redis}
}

func (r *RefreshTokenRepo) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, device, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(
		ctx, query,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.Device,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
	)

	return err
}

func (r *RefreshTokenRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, device, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > now()
	`

	var refreshToken model.RefreshToken

	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.Device,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &refreshToken, err
}

func (r *RefreshTokenRepo) IsBlacklisted(ctx context.Context, tokenHash string) (bool, error) {
	exists, err := r.redis.Exists(ctx, "blacklist:"+tokenHash).Result()

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *RefreshTokenRepo) RevokeRefreshToken(ctx context.Context, tokenHash string, expiresAt time.Time) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token_hash = $1
	`

	_, err := r.pool.Exec(ctx, query, tokenHash)

	if err != nil {
		return err
	}

	ttl := time.Until(expiresAt)

	if ttl > 0 {
		err = r.redis.Set(ctx, "blacklist:"+tokenHash, "1", ttl).Err()

		if err != nil {
			return err
		}
	}

	return nil
}
