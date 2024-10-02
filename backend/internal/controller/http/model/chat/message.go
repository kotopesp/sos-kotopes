package chat

import "time"

type (
	Message struct {
		// ID        int       `json:"id"`
		UserID  int    `json:"user_id"`
		ChatID  int    `json:"chat_id"`
		Content string `json:"message_content"`
		// IsDeleted bool      `json:"is_deleted"`
		// DeletedAt time.Time `json:"deleted_at"`
		CreatedAt  time.Time `json:"created_at"`
		SenderName string    `json:"sender_name"`
		// UpdatedAt time.Time `json:"updated_at"`
		// IsRead bool `json:"is_read"`
	}

	MessageRead struct {
		ID        int       `json:"id"`
		MessageID int       `json:"message_id"`
		UserID    int       `json:"user_id"`
		ReadAt    time.Time `json:"read_at"`
	}
)
