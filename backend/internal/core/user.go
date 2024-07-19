package core

import (
	"context"
	"time"
)

type (
	User struct {
		ID           int       `gorm:"column:id"`
		ExternalID   *int      `gorm:"column:ext_id"`
		Username     string    `gorm:"column:username"`
		Firstname    string    `gorm:"column:firstname"`
		Lastname     string    `gorm:"column:lastname"`
		Photo        *[]byte   `gorm:"column:photo"`
		PasswordHash string    `gorm:"column:password_hash"`
		Description  string    `gorm:"column:description"`
		IsDeleted    bool      `gorm:"is_deleted"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}

	UserStore interface {
		GetUserByID(ctx context.Context, id int) (data User, err error)
		GetUserByUsername(ctx context.Context, username string) (data User, err error)
		GetUserByExternalID(ctx context.Context, extID int) (data User, err error)
		AddUser(ctx context.Context, user User) (id int, err error)
	}
)

// table name in db for gorm
func (User) TableName() string {
	return "users"
}
