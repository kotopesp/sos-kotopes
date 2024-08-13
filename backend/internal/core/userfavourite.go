package core

import "context"

// todo
type (
	FavouriteUser struct {
		ID        int    `gorm:"column:id"`
		PersonID  int    `gorm:"column:person_id"`
		UserID    int    `gorm:"column:user_id"`
		createdAt string `gorm:"column:created_at"`
	}
	UserFavouriteStore interface {
		AddUserToFavourite(ctx context.Context, favouriteUserID int, userID int) (favouriteUser User, err error)
		GetFavouriteUsers(ctx context.Context, userID int) (favouriteUsers []User, err error)
		DeleteUserFromFavourite(ctx context.Context, favouriteUserID int, userID int) (err error)
	}
	UserFavouriteService interface {
		AddUserToFavourite(ctx context.Context, favouriteUserID int, userID int) (favouriteUser User, err error)
		GetFavouriteUsers(ctx context.Context, userID int) (favouriteUsers []User, err error)
		DeleteUserFromFavourite(ctx context.Context, favouriteUserID int, userID int) (err error)
	}
	GetFavourites struct {
		Count  *int
		Offset *int
		Sort   *string
	}
)

func (FavouriteUser) TableName() string {
	return "favourite_persons"
}
