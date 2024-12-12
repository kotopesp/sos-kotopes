package vet_review

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

func New(pg *postgres.Postgres) core.VetReviewsStore {
	return &store{pg}
}

func (s *store) CreateReview(ctx context.Context, review core.VetReviews) error {
	if err := s.DB.WithContext(ctx).Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) SoftDeleteReviewByID(ctx context.Context, id int) error {
	err := s.DB.WithContext(ctx).Model(&core.VetReviews{}).Where("id = ?", id).Updates(core.VetReviews{IsDeleted: true, DeletedAt: time.Now()}).Error
	return err
}

func (s *store) UpdateReviewByID(ctx context.Context, review core.UpdateVetReviews) (core.VetReviews, error) {
	review.UpdatedAt = time.Now()

	var updatedVetReview core.VetReviews

	result := s.DB.WithContext(ctx).Model(&core.VetReviews{}).Where("id = ? AND is_deleted = ?", review.ID, false).Updates(review).First(&updatedVetReview, review.ID)
	if result.Error != nil {
		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.VetReviews{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, result.Error.Error())
		return core.VetReviews{}, result.Error
	}
	return updatedVetReview, nil
}

func (s *store) GetAllReviews(ctx context.Context, params core.GetAllVetReviewsParams) ([]core.VetReviews, error) {
	var reviews []core.VetReviews

	query := s.DB.WithContext(ctx).Model(&core.VetReviews{}).Where("is_deleted = ?", false)

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}

	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	err := query.Find(&reviews).Error
	return reviews, err
}

func (s *store) GetByIDReview(ctx context.Context, id int) (core.VetReviews, error) {
	var review = core.VetReviews{ID: id}

	if err := s.DB.WithContext(ctx).First(&review).Error; err != nil {
		return core.VetReviews{}, err
	}

	return review, nil
}
