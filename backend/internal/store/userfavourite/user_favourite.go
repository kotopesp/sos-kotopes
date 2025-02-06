package userfavourite

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserFavouriteStore {
	return &store{pg}
}

func (s *store) AddUserToFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}

func (s *store) GetFavouriteUsers(ctx context.Context, userID int, params core.GetFavourites) (favouriteUsers []core.User, err error) {
	panic("implement me")
}

func (s *store) DeleteUserFromFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}
