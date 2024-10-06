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
	}

	MessageRead struct {
		ID        int       `gorm:"column:id;primaryKey"`
		MessageID int       `gorm:"column:message_id"`
		UserID    int       `gorm:"column:user_id"`
		ReadAt    time.Time `gorm:"column:read_at"`
	}

	ChatMember struct {
		UserID    int       `gorm:"column:user_id;primaryKey"`
		ChatID    int       `gorm:"column:chat_id;primaryKey"`
		IsDeleted bool      `gorm:"column:is_deleted"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	ChatStore interface {
		GetAllChats(ctx context.Context, sortType string, userID int) (chats []chat.Chat, err error)
		GetChatWithUsersByID(ctx context.Context, id int) (chat chat.Chat, err error)
		CreateChat(ctx context.Context, data chat.Chat) (chat chat.Chat, err error)
		FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error)
		DeleteChat(ctx context.Context, id int) (err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
	}

	ChatService interface {
		GetAllChats(ctx context.Context, sortType string, userID int) (chats []chat.Chat, total int, err error)
		GetChatWithUsersByID(ctx context.Context, id int) (chat chat.Chat, err error)
		CreateChat(ctx context.Context, data chat.Chat, userIds []int) (chat chat.Chat, err error)
		FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error)
		DeleteChat(ctx context.Context, id int) (err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
	}

	MessageStore interface {
		MarkMessagesAsRead(ctx context.Context, chatID, userID int) error
		GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error)
		GetAllMessages(ctx context.Context, chatID int, sortType string, searchText string) (messages []chat.Message, err error)
		CreateMessage(ctx context.Context, data chat.Message) (message chat.Message, err error)
		UpdateMessage(ctx context.Context, chatID int, messageID int, data string) (message Message, err error)
		DeleteMessage(ctx context.Context, chatID int, messageID int) (err error)
	}

	MessageService interface {
		MarkMessagesAsRead(ctx context.Context, chatID, userID int) error
		GetUnreadMessageCount(ctx context.Context, chatID, userID int) (int64, error)
		GetAllMessages(ctx context.Context, chatID int, sortType string, searchText string) (messages []chat.Message, total int, err error)
		CreateMessage(ctx context.Context, data chat.Message) (message chat.Message, err error)
		UpdateMessage(ctx context.Context, chatID int, messageID int, data string) (message Message, err error)
		DeleteMessage(ctx context.Context, chatID int, messageID int) (err error)
	}

	ChatMemberStore interface {
		GetAllMembers(ctx context.Context, chatID int) (members []ChatMember, err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
		UpdateMemberInfo(ctx context.Context, chatID int, userID int) (member ChatMember, err error)
		DeleteMemberFromChat(ctx context.Context, chatID int, userID int) (err error)
	}

	ChatMemberService interface {
		GetAllMembers(ctx context.Context, chatID int) (members []ChatMember, total int, err error)
		AddMemberToChat(ctx context.Context, data ChatMember) (member ChatMember, err error)
		UpdateMemberInfo(ctx context.Context, chatID int, userID int) (member ChatMember, err error)
		DeleteMemberFromChat(ctx context.Context, chatID int, userID int) (err error)
	}
)

func (Chat) TableName() string {
	return "chats"
}

func (Message) TableName() string {
	return "messages"
}

func (ChatMember) TableName() string {
	return "chat_members"
}
