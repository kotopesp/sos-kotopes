package chat

import "time"

type (
	Message struct {
		UserID      int       `json:"user_id" form:"user_id"`
		ChatID      int       `json:"chat_id" form:"chat_id"`
		TextContent string    `json:"message_content"`
		CreatedAt   time.Time `json:"created_at" form:"created_at"`
		SenderName  string    `json:"sender_name" form:"sender_name"`
		IsAudio     bool      `json:"is_audio" form:"is_audio"`
		AudioBytes  []byte    `json:"audio_bytes" form:"audio_bytes"`
	}

	MessageRead struct {
		ID        int       `json:"id"`
		MessageID int       `json:"message_id"`
		UserID    int       `json:"user_id"`
		ReadAt    time.Time `json:"read_at"`
	}
)
