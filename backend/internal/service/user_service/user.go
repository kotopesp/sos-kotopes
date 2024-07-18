package user_service

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type Service struct {
	userStore core.UserStore
}

func (s *Service) UpdateUser(ctx context.Context, id int, update user.UpdateUser) error {
	return s.userStore.UpdateUser(ctx, id, update)
}

func NewUserService(store core.UserStore) core.UserService {
	return &Service{
		userStore: store,
	}
}
