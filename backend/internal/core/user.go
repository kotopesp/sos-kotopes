package core

import "context"

type (
	User struct {
		Id           int    `gorm:"primary key;autoIncrement" db:"id"`
		Username     string `db:"username"`
		PasswordHash string `db:"password_hash"`
		CreatedAt    string `db:"created_at"`
	}
	UserStore interface {
		ChangeName(ctx context.Context, id int, name string) (err error)
	}
	UserService interface {
		ChangeName(ctx context.Context, id int, name string) (err error)
	}
)
