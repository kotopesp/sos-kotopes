package user

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type Service struct {
	userStore core.UserStore
}

func (s *Service) UpdateUser(ctx context.Context, id int, update user.UpdateUser) (err error) {
	return s.userStore.UpdateUser(ctx, id, update)
}
func (s *Service) GetUser(ctx context.Context, id int) (user core.User, err error) {
	return s.userStore.GetUser(ctx, id)
}
func (s *Service) GetUserPosts(ctx context.Context, id int) (posts []core.Post, err error) {
	return s.userStore.GetUserPosts(ctx, id)
}

func NewUserService(store core.UserStore) core.UserService {
	return &Service{
		userStore: store,
	}
}
