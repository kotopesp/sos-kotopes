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

func New(pg *postgres.Postgres) core.PostStore {
	return &store{pg}
}

func (s *store) GetAll(ctx context.Context, params core.GetAllPostsParams) ([]core.Post, error) {
	var posts []core.Post
	query := s.DB.WithContext(ctx).Model(&core.Post{})

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}

	if params.SearchTerm != nil {
		query = query.Where("title ILIKE ?", "%" + *params.SearchTerm + "%")
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Post, error) {
	var post core.Post = core.Post{ID : id}
	
	if err := s.DB.WithContext(ctx).First(&post).Error; err != nil {
		return core.Post{}, err
	}
	return post, nil
}

func (s *store) Create(ctx context.Context, post core.Post) (core.Post, error) {
    if err := s.DB.WithContext(ctx).Create(&post).Error; err != nil {
        return post, err
    }
    return post, nil
}

func (s *store) Update(ctx context.Context, post core.Post) (core.Post, error) {
    if err := s.DB.WithContext(ctx).Save(&post).Error; err != nil {
        return post, err
    }
    return post, nil
}

func (s *store) Delete(ctx context.Context, id int) error {
    if err := s.DB.WithContext(ctx).Delete(&core.Post{}, id).Error; err != nil {
        return err
    }
    return nil
}