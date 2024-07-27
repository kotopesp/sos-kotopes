package user_favourite_service

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type Service struct {
	userStore core.UserFavouriteStore
}

func New(store core.UserFavouriteStore) core.UserFavouriteService {
	return &Service{
		userStore: store,
	}
}

func (s *Service) AddUserToFavourite(ctx context.Context, personId int, userId int) (err error) {
	return nil
}
func (s *Service) GetFavouriteUsers(ctx context.Context, userId int) (persons []core.FavouriteUser, err error) {
	return nil, nil
}
func (s *Service) DeleteUserFromFavourite(ctx context.Context, personId int, userId int) (err error) {
	return nil
}
