package userfavourite

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	userFavouriteStore core.UserFavouriteStore
}

func New(store core.UserFavouriteStore) core.UserFavouriteService {
	return &service{
		userFavouriteStore: store,
	}
}

func (s *service) AddUserToFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}

func (s *service) GetFavouriteUsers(ctx context.Context, userID int, params core.GetFavourites) (favouriteUsers []core.User, err error) {
	panic("implement me")
}
func (s *service) DeleteUserFromFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}
