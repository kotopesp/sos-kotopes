package core

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
)

type (
	User struct {
		Id           int    `gorm:"primary key;autoIncrement" db:"id"`
		Username     string `db:"username"`
		PasswordHash string `db:"password_hash"`
		CreatedAt    string `db:"created_at"`
	}
	UserStore interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
	}
	UserService interface {
		UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error)
	}
)
