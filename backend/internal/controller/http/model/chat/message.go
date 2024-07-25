package chat

import "time"

type (
	Message struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		ChatID    int       `json:"chat_id"`
		Content   string    `json:"content"`
		IsDeleted bool      `json:"is_deleted"`
		DeletedAt time.Time `json:"deleted_at"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
