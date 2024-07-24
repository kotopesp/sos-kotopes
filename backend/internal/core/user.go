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
	FavouriteUser struct {
		Id        int    `gorm:"primary key;autoIncrement" db:"id"`
		personId  int    `db:"person_id"`
		userId    int    `db:"user_id"`
		createdAt string `db:"created_at"`
	}
	FavouriteUserStore interface {
		AddPersonToFavourite(ctx context.Context, personId int, userId int) (err error)
		GetFavouritePersons(ctx context.Context, userId int) (persons []FavouriteUser, err error)
		DeletePersonFromFavourite(ctx context.Context, personId int, userId int) (err error)
	}
	FavouriteUserService interface {
		AddPersonToFavourite(ctx context.Context, personId int, userId int) (err error)
		GetFavouritePersons(ctx context.Context, userId int) (persons []FavouriteUser, err error)
		DeletePersonFromFavourite(ctx context.Context, personId int, userId int) (err error)
	}
)
