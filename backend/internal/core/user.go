package core

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
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
	UserStore interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
		GetUser(ctx context.Context, id int) (user User, err error)
		GetUserPosts(ctx context.Context, id int) (posts []Post, err error)
		GetUserByID(ctx context.Context, id int) (data User, err error)
		GetUserByUsername(ctx context.Context, username string) (data User, err error)
		GetUserByExternalID(ctx context.Context, externalID int) (data ExternalUser, err error)
		AddUser(ctx context.Context, user User) (userID int, err error)
		AddExternalUser(ctx context.Context, user User, externalUserID int, authProvider string) (userID int, err error)
	}
	UserService interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
		GetUser(ctx context.Context, id int) (user User, err error)
		GetUserPosts(ctx context.Context, id int) (posts []Post, err error)
	}
	FavouriteUser struct {
		Id        int    `gorm:"primary key;autoIncrement" db:"id"`
		personId  int    `db:"person_id"`
		userId    int    `db:"user_id"`
		createdAt string `db:"created_at"`
	}
	UserFavouriteStore interface {
		AddUserToFavourite(ctx context.Context, personId int, userId int) (err error)
		GetFavouriteUsers(ctx context.Context, userId int, params user.GetFavourites) (favouriteUsers []user.User, err error)
		DeleteUserFromFavourite(ctx context.Context, personId int, userId int) (err error)
	}
	UserFavouriteService interface {
		AddUserToFavourite(ctx context.Context, personId int, userId int) (err error)
		GetFavouriteUsers(ctx context.Context, userId int, params user.GetFavourites) (favouriteUsers []user.User, err error)
		DeleteUserFromFavourite(ctx context.Context, personId int, userId int) (err error)
	}
	GetFavourites struct {
		Count  *int
		Offset *int
		Sort   *string
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
)

// TableName table name in db for gorm
func (User) TableName() string {
	return "users"
}

func (ExternalUser) TableName() string {
	return "external_users"
}
