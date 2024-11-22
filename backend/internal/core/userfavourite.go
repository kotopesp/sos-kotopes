package core

import "context"

// todo
type (
	FavouriteUser struct {
		ID        int    `gorm:"column:id"`
		personID  int    `gorm:"column:person_id"`
		userID    int    `gorm:"column:user_id"`
		createdAt string `gorm:"column:created_at"`
	}
	UserFavouriteStore interface {
		AddUserToFavourite(ctx context.Context, personID int, userID int) (err error)
		GetFavouriteUsers(ctx context.Context, userID int, params GetFavourites) (favouriteUsers []User, err error)
		DeleteUserFromFavourite(ctx context.Context, personID int, userID int) (err error)
	}
	UserFavouriteService interface {
		AddUserToFavourite(ctx context.Context, personID int, userID int) (err error)
		GetFavouriteUsers(ctx context.Context, userID int, params GetFavourites) (favouriteUsers []User, err error)
		DeleteUserFromFavourite(ctx context.Context, personID int, userID int) (err error)
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
