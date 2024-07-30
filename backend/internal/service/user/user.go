package user

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	userStore core.UserStore
}

func New(store core.UserStore) core.UserService {
	return &service{
		userStore: store,
	}
}

func (s *service) UpdateUser(ctx context.Context, id int, update core.UpdateUser) (updatedUser core.User, err error) {
	return s.userStore.UpdateUser(ctx, id, update)
}

func (s *service) GetUser(ctx context.Context, id int) (user core.User, err error) {
	return s.userStore.GetUser(ctx, id)
}

func (s *service) GetUserPosts(ctx context.Context, id int) (postsDetails []core.PostDetails, err error) {
	posts, err := s.userStore.GetUserPosts(ctx, id)
	for _, post := range posts {
		user, errLoop := s.userStore.GetUserByID(ctx, post.UserID)
		err = errLoop
		userName := user.Username
		postsDetails = append(postsDetails, post.ToPostDetails(userName))
	}
	return postsDetails, err
}
