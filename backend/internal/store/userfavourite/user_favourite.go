package userfavourite

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserFavouriteStore {
	return &Store{pg}
}

func (s *Store) AddUserToFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}

func (s *Store) GetFavouriteUsers(ctx context.Context, userID int, params core.GetFavourites) (favouriteUsers []core.User, err error) {
	panic("implement me")
}

func (s *Store) DeleteUserFromFavourite(ctx context.Context, personID, userID int) (err error) {
	panic("implement me")
}
