package refreshsession

import (
	"context"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

const (
	maxSessionsPerUser = 5
)

func (s *store) CountSessionsAndDelete(ctx context.Context, userID int) error {
	tx := s.DB.WithContext(ctx).Begin()

	var sessionsCounter int64
	if err := tx.Model(&core.RefreshSession{}).
		Where("user_id=?", userID).
		Count(&sessionsCounter).Error; err != nil {
		tx.Rollback()
		return err
	}

	if sessionsCounter < maxSessionsPerUser {
		return tx.Commit().Error
	}

	if err := tx.Where("user_id=?", userID).Delete(&core.RefreshSession{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *store) UpdateRefreshSession(
	ctx context.Context,
	param core.UpdateRefreshSessionParam,
	rs core.RefreshSession,
) error {
	tx := s.DB.WithContext(ctx).Begin()

	err := param(tx)
	if err != nil {
		return err
	}

	if err := tx.Create(&rs).Error; err != nil {
		tx.Rollback()
		return err
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
