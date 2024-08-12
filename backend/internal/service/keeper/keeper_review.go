package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (s *service) GetAllReviews(ctx context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviewsDetails, error) {
	keeperReviews, err := s.keeperReviewsStore.GetAllReviews(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	keeperReviewsDetails := make([]core.KeeperReviewsDetails, len(keeperReviews))

	for i, review := range keeperReviews {
		keeperReviewUser, err := s.userStore.GetUserByID(ctx, review.AuthorID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		keeperReviewsDetails[i] = core.KeeperReviewsDetails{
			Review: review,
			User:   keeperReviewUser,
		}
	}

	return keeperReviewsDetails, nil
}

func (s *service) CreateReview(ctx context.Context, review core.KeeperReviews) error {
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}

	if review.Grade < 1 || review.Grade > 5 {
		return core.ErrReviewGradeBounds
	}

	return s.keeperReviewsStore.CreateReview(ctx, review)
}

func (s *service) DeleteReviewByID(ctx context.Context, id int) error {
	return s.keeperReviewsStore.DeleteReviewByID(ctx, id)
}

func (s *service) SoftDeleteReviewByID(ctx context.Context, id int) error {
	return s.keeperReviewsStore.SoftDeleteReviewByID(ctx, id)
}

func (s *service) UpdateReviewByID(ctx context.Context, review core.UpdateKeeperReviews) (core.KeeperReviewsDetails, error) {
	updatedKeeperReview, err := s.keeperReviewsStore.UpdateReviewByID(ctx, review)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeeperReviewsDetails{}, err
	}

	keeperReviewUser, err := s.userStore.GetUserByID(ctx, updatedKeeperReview.AuthorID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeeperReviewsDetails{}, err
	}

	return core.KeeperReviewsDetails{
		Review: updatedKeeperReview,
		User:   keeperReviewUser,
	}, nil
}
