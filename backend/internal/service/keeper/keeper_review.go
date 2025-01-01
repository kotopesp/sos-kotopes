package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *service) GetAllReviews(ctx context.Context, keeperID int, params core.GetAllKeeperReviewsParams) (data []core.KeeperReview, err error) {
	keeperReviews, err := s.keeperReviewStore.GetAllReviews(ctx, keeperID, params)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return keeperReviews, nil
}

func (s *service) CreateReview(ctx context.Context, review core.KeeperReview) (data core.KeeperReview, err error) {
	review.CreatedAt = time.Now()

	if review.Grade < 1 || review.Grade > 5 {
		logger.Log().Debug(ctx, core.ErrReviewGradeBounds.Error())
		return core.KeeperReview{}, core.ErrReviewGradeBounds
	}

	if err := s.keeperReviewStore.CreateReview(ctx, review); err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.KeeperReview{}, err
	}

	return review, nil
}

func (s *service) DeleteReview(ctx context.Context, id, userID int) error {
	storedReview, err := s.keeperReviewStore.GetReviewByID(ctx, id)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	if storedReview.AuthorID != userID {
		logger.Log().Debug(ctx, core.ErrKeeperReviewUserIDMissmatch.Error())
		return core.ErrKeeperReviewUserIDMissmatch
	} else if storedReview.IsDeleted {
		logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
		return core.ErrRecordNotFound
	}

	return s.keeperReviewStore.DeleteReview(ctx, id)
}

func (s *service) UpdateReview(ctx context.Context, id, userID int, review core.KeeperReview) (data core.KeeperReview, err error) {
	storedReview, err := s.keeperReviewStore.GetReviewByID(ctx, id)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.KeeperReview{}, err
	}

	if storedReview.AuthorID != userID {
		logger.Log().Debug(ctx, core.ErrKeeperReviewUserIDMissmatch.Error())
		return core.KeeperReview{}, core.ErrKeeperReviewUserIDMissmatch
	} else if storedReview.IsDeleted {
		logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
		return core.KeeperReview{}, core.ErrRecordNotFound
	}

	updatedReview, err := s.keeperReviewStore.UpdateReview(ctx, id, review)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.KeeperReview{}, err
	}

	return updatedReview, nil
}
