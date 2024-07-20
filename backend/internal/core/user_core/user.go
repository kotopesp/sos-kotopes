package user_core

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core/post_core"
)

type (
	User struct {
		Id        int    `gorm:"primary key;autoIncrement" db:"id"`
		Username  string `db:"username"`
		Password  string `db:"password"`
		CreatedAt string `db:"created_at"`
	}
	UserStore interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
		GetUser(ctx context.Context, id int) (user User, err error)
		GetUserPosts(ctx context.Context, id int) ([]post_core.Post, error)
	}
	UserService interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
		GetUser(ctx context.Context, id int) (user User, err error)
		GetUserPosts(ctx context.Context, id int) ([]post_core.Post, error)
	}
)
