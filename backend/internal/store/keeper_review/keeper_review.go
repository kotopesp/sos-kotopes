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

func New(pg *postgres.Postgres) core.KeeperReviewStore {
	return &store{pg}
}

func (s *store) CreateReview(ctx context.Context, review core.KeeperReview) error {
	if err := s.DB.WithContext(ctx).Create(&review).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteReview(ctx context.Context, id int) error {
	now := time.Now()
	err := s.DB.WithContext(ctx).Model(&core.KeeperReview{}).Where("id = ?", id).Updates(core.KeeperReview{IsDeleted: true, DeletedAt: &now}).Error

	return err
}

func (s *store) UpdateReview(ctx context.Context, id int, review core.UpdateKeeperReview) (data core.KeeperReview, err error) {
	var updatedReview core.KeeperReview
	updatedReview.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Model(&core.KeeperReview{}).Preload("Author").Preload("User").Where("id = ? AND is_deleted = ?", id, false).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.KeeperReview{}, err
	}

	if err := s.DB.WithContext(ctx).Updates(review).First(&updatedReview).Error; err != nil {
		switch {
		case errors.Is(err, core.ErrRecordNotFound):
			logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
			return core.KeeperReview{}, core.ErrRecordNotFound
		default:
			logger.Log().Debug(ctx, err.Error())
			return core.KeeperReview{}, err
		}
	}

	return updatedReview, nil
}

func (s *store) GetAllReviews(ctx context.Context, keeperID int, params core.GetAllKeeperReviewsParams) (data []core.KeeperReview, err error) {
	var reviews []core.KeeperReview

	query := s.DB.WithContext(ctx).Model(&core.KeeperReview{}).Where("is_deleted = ? AND keeper_id", false, keeperID).Preload("Author").Preload("Keeper")

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (s *store) GetReviewByID(ctx context.Context, id int) (core.KeeperReview, error) {
	var review = core.KeeperReview{ID: id}

	if err := s.DB.WithContext(ctx).Preload("Author").Preload("Keeper").First(&review).Error; err != nil {
		return core.KeeperReview{}, err
	}

	return review, nil
}
