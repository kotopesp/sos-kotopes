package chat

import "time"

type (
	ChatMember struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		ChatID    int       `json:"chat_id"`
		IsDeleted bool      `json:"is_deleted"`
		DeletedAt time.Time `json:"deleted_at"`
		CreatedAt time.Time `json:"created_at"`
	}
)
