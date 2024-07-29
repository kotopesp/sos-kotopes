package keeperreviewstore

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.KeeperReviewsStore {
	return &store{pg}
}

func (s *store) Create(ctx *context.Context, review core.KeeperReviews) error {
	if err := s.DB.WithContext(*ctx).Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteByID(ctx *context.Context, id int) error {
	result := s.DB.WithContext(*ctx).Delete(core.KeeperReviews{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (s *store) SoftDeleteByID(ctx *context.Context, id int) error {
	err := s.DB.WithContext(*ctx).Model(&core.KeeperReviews{}).Where("id = ?", id).Updates(core.KeeperReviews{IsDeleted: true, DeletedAt: time.Now()}).Error

	return err
}

func (s *store) UpdateByID(ctx *context.Context, review core.KeeperReviews) error {
	result := s.DB.WithContext(*ctx).Model(&core.KeeperReviews{}).Where("id = ? AND is_deleted = ?", review.ID, false).Updates(review)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (s *store) GetAll(ctx *context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	var reviews []core.KeeperReviews

	query := s.DB.WithContext(*ctx).Model(&core.KeeperReviews{}).Where("is_deleted = ?", false)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	err := query.Find(&reviews).Error
	return reviews, err
}
