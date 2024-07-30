package user_favourite

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	userStore core.UserFavouriteStore
}

func New(store core.UserFavouriteStore) core.UserFavouriteService {
	return &service{
		userStore: store,
	}
}

func (s *service) AddUserToFavourite(ctx context.Context, personId int, userId int) (err error) {
	panic("implement me")
}

func (s *service) GetFavouriteUsers(ctx context.Context, userId int, params core.GetFavourites) (favouriteUsers []core.User, err error) {
	panic("implement me")
}
func (s *service) DeleteUserFromFavourite(ctx context.Context, personId int, userId int) (err error) {
	panic("implement me")
}
