package chat

import "time"

type (
	Message struct {
		UserID     int       `json:"user_id"`
		ChatID     int       `json:"chat_id"`
		Content    string    `json:"message_content"`
		CreatedAt  time.Time `json:"created_at"`
		SenderName string    `json:"sender_name"`
	}

	MessageRead struct {
		ID        int       `json:"id"`
		MessageID int       `json:"message_id"`
		UserID    int       `json:"user_id"`
		ReadAt    time.Time `json:"read_at"`
	}
)
