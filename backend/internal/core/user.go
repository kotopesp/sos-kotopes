package core

import (
	"context"
	"time"
)

type (
	User struct {
		ID          int       `gorm:"column:id"`
		Username    string    `gorm:"column:username"`
		Password    string    `gorm:"column:password"`
		Description string    `gorm:"column:description"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	UserStore interface {
		GetUserByID(ctx context.Context, id int) (data User, err error)
		GetUserByUsername(ctx context.Context, username string) (data User, err error)
	}
)

// table name in db for gorm
func (User) TableName() string {
	return "users"
}
