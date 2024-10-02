package user

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	userStore          core.UserStore
	userFavouriteStore core.UserFavouriteStore
}

func New(store core.UserStore, favouriteStore core.UserFavouriteStore) core.UserService {
	return &service{
		userStore:          store,
		userFavouriteStore: favouriteStore,
	}
}

func (s *service) GetUser(ctx context.Context, id int) (user chat.User, err error) {
	return s.userStore.GetUser(ctx, id)
}

func (s *service) UpdateUser(ctx context.Context, id int, update core.UpdateUser) (updatedUser chat.User, err error) {
	return s.userStore.UpdateUser(ctx, id, update)
}
