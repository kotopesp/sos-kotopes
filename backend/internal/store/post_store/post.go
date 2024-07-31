package poststore

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.PostStore {
	return &store{pg}
}

func (s *store) GetPostByID(ctx context.Context, id int) (data core.Post, err error) {
	var post core.Post

	if err := s.DB.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Post{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Post{}, err
	}

	return post, nil
}
