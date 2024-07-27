package user_favourite_store

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type Store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserFavouriteStore {
	return &Store{pg}
}
func (r *Store) AddUserToFavourite(ctx context.Context, personId int, userId int) (err error) {
	return nil
}
func (r *Store) GetFavouriteUsers(ctx context.Context, userId int) (persons []core.FavouriteUser, err error) {
	return nil, nil
}
func (r *Store) DeleteUserFromFavourite(ctx context.Context, personId int, userId int) (err error) {
	return nil
}
