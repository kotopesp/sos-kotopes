package core

import (
	"context"
	"time"
)

type (
	User struct {
		ID           int        `gorm:"column:id"`
		Username     string     `gorm:"column:username"`
		Firstname    *string    `gorm:"column:firstname"`
		Lastname     *string    `gorm:"column:lastname"`
		Photo        *[]byte    `gorm:"column:photo"`
		PasswordHash string     `gorm:"column:password_hash"`
		Description  *string    `gorm:"column:description"`
		IsDeleted    bool       `gorm:"is_deleted"`
		CreatedAt    time.Time  `gorm:"column:created_at"`
		UpdatedAt    time.Time  `gorm:"column:updated_at"`
		DeletedAt    *time.Time `gorm:"column:deleted_at"`
	}

	ExternalUser struct {
		ID           int    `gorm:"id"`
		UserID       int    `gorm:"user_id"`
		ExternalID   int    `gorm:"external_id"`
		AuthProvider string `gorm:"auth_provider"`
	}

	UserStore interface {
		GetUserByID(ctx context.Context, id int) (data User, err error)
		GetUserByUsername(ctx context.Context, username string) (data User, err error)
		GetUserByExternalID(ctx context.Context, externalID int) (data ExternalUser, err error)
		AddUser(ctx context.Context, user User) (userID int, err error)
		AddExternalUser(ctx context.Context, user User, externalUserID int, authProvider string) (userID int, err error)
	}
)

// TableName table name in db for gorm
func (User) TableName() string {
	return "users"
}

func (ExternalUser) TableName() string {
	return "external_users"
}
