package conversation

import "time"

type (
	Message struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		IsDeleted bool      `json:"is_deleted"`
		ChatID    int       `json:"chat_id"`
	}
)
