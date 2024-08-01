package keeperreviewservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	KeeperReviewsStore core.KeeperReviewsStore
}

func New(keeperReviewsStore core.KeeperReviewsStore) core.KeeperReviewsStore {
	return &service{KeeperReviewsStore: keeperReviewsStore}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	return s.KeeperReviewsStore.GetAll(ctx, params)
}

func (s *service) Create(ctx context.Context, review core.KeeperReviews) error {
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}

	if review.Grade < 1 || review.Grade > 5 {
		return core.ErrReviewGradeBounds
	}

	return s.KeeperReviewsStore.Create(ctx, review)
}

func (s *service) DeleteByID(ctx context.Context, id int) error {
	return s.KeeperReviewsStore.DeleteByID(ctx, id)
}

func (s *service) SoftDeleteByID(ctx context.Context, id int) error {
	return s.KeeperReviewsStore.SoftDeleteByID(ctx, id)
}

func (s *service) UpdateByID(ctx context.Context, review core.KeeperReviews) error {
	if review.Grade < 1 || review.Grade > 5 {
		return core.ErrReviewGradeBounds
	}

	return s.KeeperReviewsStore.UpdateByID(ctx, review)
}