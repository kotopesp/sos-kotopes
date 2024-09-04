package refreshsession

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

// CreateRefreshSession implements core.RefreshSessionStore.
func (s *store) CreateRefreshSession(ctx context.Context, rs core.RefreshSession) error {
	if err := s.DB.WithContext(ctx).Create(&rs).Error; err != nil {
		return err
	}

	return nil
}

func (s *store) UpdateRefreshSession(ctx context.Context, oldSessionID int, rs core.RefreshSession) error {
	tx := s.DB.WithContext(ctx).Begin()

	if err := tx.Model(&core.RefreshSession{}).Delete("id=?", oldSessionID).Error; err != nil {
		return tx.Rollback().Error
	}

	if err := tx.Create(&rs).Error; err != nil {
		return tx.Rollback().Error
	}

	return tx.Commit().Error
}

func (s *store) GetRefreshSessionByToken(ctx context.Context, token string) (data core.RefreshSession, err error) {
	if err := s.DB.WithContext(ctx).First(&data, "refresh_token=?", token).Error; err != nil {
		return data, err
	}

	return data, nil
}

func New(pg *postgres.Postgres) core.RefreshSessionStore {
	return &store{pg}
}
