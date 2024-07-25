package keeperreviewStore

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/postgres"
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

func (s *store) DeleteById(ctx *context.Context, id int) error {
	result := s.DB.WithContext(*ctx).Delete(core.KeeperReviews{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (s *store) UpdateById(ctx *context.Context, review core.KeeperReviews) error {
	result := s.DB.WithContext(*ctx).Model(&core.KeeperReviews{}).Where("id = ?", review.ID).Updates(review)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

func (s *store) GetAll(ctx *context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	var reviews []core.KeeperReviews

	query := s.DB.WithContext(*ctx).Model(&core.KeeperReviews{})

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	err := query.Find(&reviews).Error
	return reviews, err
}
