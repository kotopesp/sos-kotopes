package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

func (s *service) GetAllReviews(ctx context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	return s.keeperReviewsStore.GetAllReviews(ctx, params)
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

func (s *service) UpdateReviewByID(ctx context.Context, review core.KeeperReviews) error {
	return s.keeperReviewsStore.UpdateReviewByID(ctx, review)
}
