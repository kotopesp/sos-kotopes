package core

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
)

type (
	Message struct {
		ID         int       `gorm:"column:id;primaryKey"`
		UserID     int       `gorm:"column:user_id"`
		ChatID     int       `gorm:"column:chat_id"`
		Content    string    `gorm:"column:content"`
		IsDeleted  bool      `gorm:"column:is_deleted"`
		DeletedAt  time.Time `gorm:"column:deleted_at"`
		CreatedAt  time.Time `gorm:"column:created_at"`
		UpdatedAt  time.Time `gorm:"column:updated_at"`
		SenderName string    `gorm:"column:sender_name"`
		IsAudio    bool      `gorm:"column:is_audio"`
		AudioBytes []byte    `gorm:"column:audio_bytes"`
	}

	MessageRead struct {
		ID        int       `gorm:"column:id;primaryKey"`
		MessageID int       `gorm:"column:message_id"`
		UserID    int       `gorm:"column:user_id"`
		ReadAt    time.Time `gorm:"column:read_at"`
	}

	MessageStore interface {
		MarkMessagesAsRead(ctx context.Context, chatID, userID int) error
		GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error)
		GetAllMessages(ctx context.Context, chatID int, userID int, sortType string, searchText string) (messages []chat.Message, err error)
		CreateMessage(ctx context.Context, data chat.Message) (message chat.Message, err error)
		UpdateMessage(ctx context.Context, chatID, userID, messageID int, data string) (message Message, err error)
		DeleteMessage(ctx context.Context, chatID, userID, messageID int) (err error)
	}

	MessageService interface {
		MarkMessagesAsRead(ctx context.Context, chatID, userID int) error
		GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error)
		GetAllMessages(ctx context.Context, chatID, userID int, sortType string, searchText string) (messages []chat.Message, total int, err error)
		CreateMessage(ctx context.Context, data chat.Message) (message chat.Message, err error)
		UpdateMessage(ctx context.Context, chatID, userID, messageID int, data string) (message Message, err error)
		DeleteMessage(ctx context.Context, chatID, userID, messageID int) (err error)
	}
)

func (Message) TableName() string {
	return "messages"
}
