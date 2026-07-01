package model

import "time"

type ChatInfo struct {
	ID        string    `json:"id"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
}
