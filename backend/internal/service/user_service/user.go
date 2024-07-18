package user_service

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core/user_core"
)

type Service struct {
	userStore user_core.UserStore
}

func (s *Service) UpdateUser(ctx context.Context, id int, update user.UpdateUser) error {
	return s.userStore.UpdateUser(ctx, id, update)
}
func (s *Service) GetUser(ctx context.Context, id int) (user_core.User, error) {
	return s.userStore.GetUser(ctx, id)
}

func NewUserService(store user_core.UserStore) user_core.UserService {
	return &Service{
		userStore: store,
	}
}
