package user

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type UserService struct {
	userStore core.UserStore
}

func (s *UserService) ChangeName(ctx context.Context, id int, name string) (err error) {
	return s.userStore.ChangeName(ctx, id, name)
}

func NewUserService(store core.UserStore) core.UserService {
	return &UserService{
		userStore: store,
	}
}
