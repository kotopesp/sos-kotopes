package postresponsestore

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.PostResponseStore {
	return &store{pg}
}

func (s *store) CreatePostResponse(ctx context.Context, response core.PostResponse) (core.PostResponse, error) {
	if err := s.DB.WithContext(ctx).Create(&response).Error; err != nil {
		return response, err
	}
	return response, nil
}

func (s *store) GetResponsesByPostID(ctx context.Context, postID int) ([]core.PostResponse, error) {
	var responses []core.PostResponse

	if err := s.DB.WithContext(ctx).Where("post_id = ? AND is_deleted = ?", postID, false).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *store) UpdatePostResponse(ctx context.Context, response core.PostResponse) (core.PostResponse, error) {
	if err := s.DB.WithContext(ctx).Save(&response).Error; err != nil {
		return response, err
	}
	return response, nil
}

func (s *store) DeletePostResponse(ctx context.Context, id int) error {
	if err := s.DB.WithContext(ctx).Where("id = ?", id).Delete(&core.PostResponse{}).Error; err != nil {
		return err
	}
	return nil
}
