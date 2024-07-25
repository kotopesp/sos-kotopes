package keeperReviewService

import (
	"context"

	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
)

type service struct {
	KeeperReviewsStore core.KeeperReviewsStore
}

func New(KeeperReviewsStore core.KeeperReviewsStore) core.KeeperReviewsStore {
	return &service{KeeperReviewsStore: KeeperReviewsStore}
}

func (s *service) GetAll(ctx *context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	return s.KeeperReviewsStore.GetAll(ctx, params)
}

func (s *service) Create(ctx *context.Context, review core.KeeperReviews) error {
	return s.KeeperReviewsStore.Create(ctx, review)
}

func (s *service) DeleteById(ctx *context.Context, id int) error {
	return s.KeeperReviewsStore.DeleteById(ctx, id)
}

func (s *service) UpdateById(ctx *context.Context, review core.KeeperReviews) error {
	return s.KeeperReviewsStore.UpdateById(ctx, review)
}
