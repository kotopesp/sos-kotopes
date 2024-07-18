package post

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type (
	service struct {
		PostStore core.PostStore
	}
)

func New(store core.PostStore) core.PostService {
	return &service{
		PostStore: store,
	}
}

func (s *service) GetByID(ctx context.Context, id int) (post core.Post, err error) {
	post, err = s.PostStore.GetByID(ctx, id)
	return
}

func (s *service) GetAll(ctx context.Context, params core.GetAllPostsParams) (posts []core.Post, total int, err error) {
	posts, err = s.PostStore.GetAll(ctx, params)
	if err != nil {
		return
	}
	total = len(posts)
	return
}

func (s *service) Create(ctx context.Context, post core.Post) (createdPost core.Post, err error) {
	createdPost, err = s.PostStore.Create(ctx, post)
	return
}

func (s *service) Update(ctx context.Context, post core.Post) (updatedPost core.Post, err error) {
	updatedPost, err = s.PostStore.Update(ctx, post)
	return
}

func (s *service) Delete(ctx context.Context, id int) (err error) {
	err = s.PostStore.Delete(ctx, id)
	return
}