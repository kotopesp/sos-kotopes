package core

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
)

type (
	User struct {
		Id           int    `gorm:"primary key;autoIncrement" db:"id"`
		Username     string `db:"username"`
		FirstName    string `db:"firstname"`
		LastName     string `db:"lastname"`
		Description  string `db:"description"`
		Photo        string `db:"photo"`
		PasswordHash string `db:"password"`
		IsDeleted    bool   `db:"is_deleted"`
		DeletedAt    string `db:"deleted_at"`
		CreatedAt    string `db:"created_at"`
		UpdatedAt    string `db:"updated_at"`
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
)
