package core

import (
	"context"
	"time"
)

type (
	ChatMember struct {
		UserID    int       `gorm:"column:user_id;primaryKey"`
		ChatID    int       `gorm:"column:chat_id;primaryKey"`
		IsDeleted bool      `gorm:"column:is_deleted"`
		DeletedAt time.Time `gorm:"column:deleted_at"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	ChatMemberStore interface {
		GetAllMembers(ctx context.Context, chatID, userID int) (members []ChatMember, err error)
		AddMemberToChat(ctx context.Context, data ChatMember, userID int) (member ChatMember, err error)
		UpdateMemberInfo(ctx context.Context, chatID int, userID, currentUserID int) (member ChatMember, err error)
		DeleteMemberFromChat(ctx context.Context, chatID int, userID, currentUserID int) (err error)
	}

	ChatMemberService interface {
		GetAllMembers(ctx context.Context, chatID, userID int) (members []ChatMember, total int, err error)
		AddMemberToChat(ctx context.Context, data ChatMember, userID int) (member ChatMember, err error)
		UpdateMemberInfo(ctx context.Context, chatID int, userID, currentUserID int) (member ChatMember, err error)
		DeleteMemberFromChat(ctx context.Context, chatID int, userID, currentUserID int) (err error)
	}
)

func (ChatMember) TableName() string {
	return "chat_members"
}
