package user

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (s *service) AddUserToFavourite(ctx context.Context, favouriteUserID, userID int) (user core.User, err error) {
	return s.userFavouriteStore.AddUserToFavourite(ctx, favouriteUserID, userID)
}

func (s *service) GetFavouriteUsers(ctx context.Context, userID int) (favouriteUsers []core.User, err error) {
	return s.userFavouriteStore.GetFavouriteUsers(ctx, userID)
}
func (s *service) DeleteUserFromFavourite(ctx context.Context, favouriteUserID, userID int) (err error) {
	return s.userFavouriteStore.DeleteUserFromFavourite(ctx, favouriteUserID, userID)
}
