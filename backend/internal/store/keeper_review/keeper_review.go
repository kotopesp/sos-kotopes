package keeperreviewstore

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.KeeperReviewsStore {
	return &store{pg}
}

func (s *store) CreateReview(ctx context.Context, review core.KeeperReviews) error {
	if err := s.DB.WithContext(ctx).Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) SoftDeleteReviewByID(ctx context.Context, id int) error {
	err := s.DB.WithContext(ctx).Model(&core.KeeperReviews{}).Where("id = ?", id).Updates(core.KeeperReviews{IsDeleted: true, DeletedAt: time.Now()}).Error

	return err
}

func (s *store) UpdateReviewByID(ctx context.Context, review core.UpdateKeeperReviews) (core.KeeperReviews, error) {
	review.UpdatedAt = time.Now()

	var updatedKeeperReview core.KeeperReviews

	result := s.DB.WithContext(ctx).Model(&core.KeeperReviews{}).Where("id = ? AND is_deleted = ?", review.ID, false).Updates(review).First(&updatedKeeperReview, review.ID)
	if result.Error != nil {

		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.KeeperReviews{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, result.Error.Error())
		return core.KeeperReviews{}, result.Error
	}
	return updatedKeeperReview, nil
}

func (s *store) GetAllReviews(ctx context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	var reviews []core.KeeperReviews

	query := s.DB.WithContext(ctx).Model(&core.KeeperReviews{}).Where("is_deleted = ?", false)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	err := query.Find(&reviews).Error
	return reviews, err
}

func (s *store) GetByIDReview(ctx context.Context, id int) (core.KeeperReviews, error) {
	var review = core.KeeperReviews{ID: id}

	if err := s.DB.WithContext(ctx).First(&review).Error; err != nil {
		return core.KeeperReviews{}, err
	}

	return review, nil
}
