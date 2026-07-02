package repository

import (
	"chat/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	ExistsPrivateChat(ctx context.Context, userID1, userID2 string) (bool, error)
	CreatePrivateChat(ctx context.Context, userID1, userID2 string) (string, error)
	GetChatsByUserID(ctx context.Context, userID string) ([]model.ChatInfo, error)
	GetCompanionID(ctx context.Context, chatID, userID string) (string, error)
}

type ChatRepo struct {
	pool *pgxpool.Pool
}

func NewChatRepo(pool *pgxpool.Pool) *ChatRepo {
	return &ChatRepo{pool}
}

func (r *ChatRepo) ExistsPrivateChat(ctx context.Context, userID1, userID2 string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM chats c
			JOIN chat_members cm1 ON c.id = cm1.chat_id AND cm1.user_id = $1
			JOIN chat_members cm2 ON c.id = cm2.chat_id AND cm2.user_id = $2
			WHERE c.type = 1
		)
	`

	err := r.pool.QueryRow(ctx, query, userID1, userID2).Scan(&exists)

	return exists, err
}

func (r *ChatRepo) CreatePrivateChat(ctx context.Context, userID1, userID2 string) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO chats (type)
		VALUES (1)
		RETURNING id
	`

	var chatID string

	err = tx.QueryRow(ctx, query).Scan(&chatID)
	if err != nil {
		return "", err
	}

	query = `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2), ($1, $3)
	`

	_, err = tx.Exec(ctx, query, chatID, userID1, userID2)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (r *ChatRepo) GetChatsByUserID(ctx context.Context, userID string) ([]model.ChatInfo, error) {
	query := `
		SELECT c.id, c.type, c.created_at, u.id, u.username
		FROM chats c
		JOIN chat_members cm_self ON cm_self.chat_id = c.id AND cm_self.user_id = $1
		JOIN chat_members cm_other ON cm_other.chat_id = c.id AND cm_other.user_id != $1
		JOIN users u ON u.id = cm_other.user_id
		ORDER BY c.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.ChatInfo

	for rows.Next() {
		var chat model.ChatInfo

		if err := rows.Scan(&chat.ID, &chat.Type, &chat.CreatedAt, &chat.UserID, &chat.Username); err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepo) GetCompanionID(ctx context.Context, chatID, userID string) (string, error) {
	var companionID string

	query := `
		SELECT user_id
		FROM chat_members
		WHERE chat_id = $1 AND user_id != $2
	`

	err := r.pool.QueryRow(ctx, query, chatID, userID).Scan(&companionID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return companionID, nil
}
