package core

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
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
		GetUserPosts(ctx context.Context, id int) (posts []Post, err error)
	}
	UserService interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
		GetUser(ctx context.Context, id int) (user User, err error)
		GetUserPosts(ctx context.Context, id int) (posts []Post, err error)
	}
)
