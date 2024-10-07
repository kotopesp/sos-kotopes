package core

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
)

type (
	Chat struct {
		ID        int       `gorm:"column:id;primaryKey"`
		ChatType  string    `gorm:"column:chat_type"`
		IsDeleted bool      `gorm:"column:is_deleted"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
		Users     []User    `gorm:"many2many:chat_members;joinForeignKey:ChatID;joinReferences:UserID"`
	}

	ChatStore interface {
		GetAllChats(ctx context.Context, sortType string, userID int) (chats []chat.Chat, err error)
		GetChatWithUsersByID(ctx context.Context, chatID, userID int) (chat chat.Chat, err error)
		CreateChat(ctx context.Context, data chat.Chat) (chat chat.Chat, err error)
		FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error)
		DeleteChat(ctx context.Context, chatID, userID int) (err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
	}

	ChatService interface {
		GetAllChats(ctx context.Context, sortType string, userID int) (chats []chat.Chat, total int, err error)
		GetChatWithUsersByID(ctx context.Context, chatID, userID int) (chat chat.Chat, err error)
		CreateChat(ctx context.Context, data chat.Chat, userIds []int) (chat chat.Chat, err error)
		FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error)
		DeleteChat(ctx context.Context, chatID, userID int) (err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
	}
)

func (Chat) TableName() string {
	return "chats"
}
