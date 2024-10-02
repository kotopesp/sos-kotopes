package core

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
)

type (
	User struct {
		ID           int        `gorm:"column:id;primaryKey"`
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
	UpdateUser struct {
		Username     *string `gorm:"column:username"`
		Firstname    *string `gorm:"column:firstname"`
		Lastname     *string `gorm:"column:lastname"`
		Description  *string `gorm:"column:description"`
		Photo        *[]byte `gorm:"column:photo"`
		PasswordHash *string `gorm:"column:password"`
	}
	UserStore interface {
		UpdateUser(ctx context.Context, id int, update UpdateUser) (updatedUser chat.User, err error)
		GetUser(ctx context.Context, id int) (user chat.User, err error)
		GetUserByID(ctx context.Context, id int) (data User, err error)
		GetUserByUsername(ctx context.Context, username string) (data User, err error)
		GetUserByExternalID(ctx context.Context, externalID int) (data ExternalUser, err error)
		AddUser(ctx context.Context, user User) (userID int, err error)
		AddExternalUser(ctx context.Context, user User, externalUserID int, authProvider string) (userID int, err error)
	}
	UserService interface {
		UpdateUser(ctx context.Context, id int, update UpdateUser) (updatedUser chat.User, err error)
		GetUser(ctx context.Context, id int) (user chat.User, err error)
	}

	ExternalUser struct {
		ID           int    `gorm:"column:id"`
		UserID       int    `gorm:"column:user_id"`
		ExternalID   int    `gorm:"column:external_id"`
		AuthProvider string `gorm:"column:auth_provider"`
	}
)

// errors
var (
	ErrNotUniqueUsername  = errors.New("username must be unique")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoResponseFromVK   = errors.New("no response from VK")
	ErrNoSuchUser         = errors.New("user does not exist")
	ErrEmptyUpdateRequest = errors.New("empty update request")
)

// TableName table name in db for gorm
func (User) TableName() string {
	return "users"
}

func (ExternalUser) TableName() string {
	return "external_users"
}
