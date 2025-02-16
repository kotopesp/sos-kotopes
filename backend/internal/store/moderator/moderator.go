package moderator

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.ModeratorStore {
	return &store{pg}
}

// GetModeratorByID - retrieves the moderator structure by their id.
func (s *store) GetModeratorByID(ctx context.Context, id int) (moderator core.Moderator, err error) {
	err = s.DB.WithContext(ctx).
		Table(moderator.TableName()).
		Where("id = ?", id).
		First(&moderator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log().Debug(ctx, err.Error())
			return moderator, core.ErrNoSuchModerator
		}
		logger.Log().Debug(ctx, err.Error())
		return moderator, err
	}

	return moderator, nil
}

// AddModerator adds moderator
func (s *store) AddModerator(ctx context.Context, userID int) (err error) {
	moderator := core.Moderator{
		UserID: userID,
	}

	err = s.DB.WithContext(ctx).
		Create(&moderator).Error

	if err != nil {
		return err
	}

	return nil
}
