package core

import (
	"context"
	"errors"
	"time"
)

type (
	Moderator struct {
		ID           int        `gorm:"column:id"`
		Username     string     `gorm:"column:username"`
		Firstname    *string    `gorm:"column:firstname"`
		Lastname     *string    `gorm:"column:lastname"`
		Photo        *[]byte    `gorm:"column:photo"`
		PasswordHash string     `gorm:"column:password_hash"`
		Description  *string    `gorm:"column:description"`
		IsDeleted    bool       `gorm:"column:is_deleted"`
		CreatedAt    time.Time  `gorm:"column:created_at"`
		UpdatedAt    time.Time  `gorm:"column:updated_at"`
		DeletedAt    *time.Time `gorm:"column:deleted_at"`
	}
	UpdateModerator struct {
		Username     *string `gorm:"column:username"`
		Firstname    *string `gorm:"column:firstname"`
		Lastname     *string `gorm:"column:lastname"`
		Description  *string `gorm:"column:description"`
		Photo        *[]byte `gorm:"column:photo"`
		PasswordHash *string `gorm:"column:password"`
		Role         *string `gorm:"column:role"`
	}
	ModeratorStore interface {
		UpdateModerator(ctx context.Context, id int, update UpdateModerator) (updatedModerator Moderator, err error)
		GetModerator(ctx context.Context, id int) (moderator Moderator, err error)
		GetModeratorByUsername(ctx context.Context, username string) (moderator Moderator, err error)
		AddModerator(ctx context.Context, moderator Moderator) (moderatorID int, err error)
	}
	ModeratorService interface {
		UpdateModerator(ctx context.Context, id int, update UpdateModerator) (updatedModerator Moderator, err error)
		GetModerator(ctx context.Context, id int) (moderator Moderator, err error)
	}
)

// Errors
var (
	ErrNoSuchModerator = errors.New("moderator does not exist")
)

// TableName table name in db for gorm
func (Moderator) TableName() string {
	return "moderators"
}
