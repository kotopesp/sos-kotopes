package post

import (
	"context"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func NewPostStore(pg *postgres.Postgres) core.PostStore {
	return &store{pg}
}

func (s *store) GetAllPosts(ctx context.Context, params core.GetAllPostsParams) ([]core.Post, int, error) {
	var posts []core.Post
	query := s.DB.WithContext(ctx).Model(&core.Post{})

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}

	if params.SearchTerm != nil {
		query = query.Where("title ILIKE ?", "%"+*params.SearchTerm+"%")
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, int(total), nil
}

func (s *store) GetPostByID(ctx context.Context, id int) (core.Post, error) {
	var post = core.Post{ID: id}

	if err := s.DB.WithContext(ctx).First(&post).Error; err != nil {
		return core.Post{}, err
	}
	return post, nil
}

func (s *store) CreatePost(ctx context.Context, post core.Post) (core.Post, error) {
	if err := s.DB.WithContext(ctx).Create(&post).Error; err != nil {
		return post, err
	}
	return post, nil
}

func (s *store) UpdatePost(ctx context.Context, post core.Post) (core.Post, error) {
	if err := s.DB.WithContext(ctx).Save(&post).Error; err != nil {
		return post, err
	}
	return post, nil
}

func (s *store) DeletePost(ctx context.Context, id int) error {
	if err := s.DB.WithContext(ctx).Where("id = ?", id).Delete(&core.Post{}).Error; err != nil {
		return err
	}
	return nil
}
