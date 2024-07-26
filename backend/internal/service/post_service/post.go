package post

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type (
	postService struct {
		PostStore core.PostStore
	}
)

func NewPostService(store core.PostStore) core.PostService {
	return &postService{
		PostStore: store,
	}
}

func (s *postService) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.Post, int, error) {
	posts, total, err := s.PostStore.GetAllPosts(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (s *postService) GetPostByID(ctx context.Context, id int) (core.Post, error) {
	return s.PostStore.GetPostByID(ctx, id)
}

func (s *postService) CreatePost(ctx context.Context, post core.Post) (core.Post, error) {
	return s.PostStore.CreatePost(ctx, post)

}

func (s *postService) UpdatePost(ctx context.Context, post core.Post) (core.Post, error) {
	return s.PostStore.UpdatePost(ctx, post)
}

func (s *postService) DeletePost(ctx context.Context, id int) error {
	return s.PostStore.DeletePost(ctx, id)
}
