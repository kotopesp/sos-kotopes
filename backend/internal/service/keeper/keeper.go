package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	keeperStore        core.KeeperStore
	KeeperReviewsStore core.KeeperReviewsService
}

func New(keeperStore core.KeeperStore, keeperReviewStore core.KeeperReviewsStore) core.KeeperService {
	return &service{
		keeperStore:        keeperStore,
		KeeperReviewsStore: keeperReviewStore,
	}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.Keepers, error) {
	if *params.SortBy == "" {
		*params.SortBy = "created_at"
	}
	if *params.SortOrder == "" {
		*params.SortBy = "desc"
	}

	return s.keeperStore.GetAll(ctx, params)
}

func (s *service) GetByID(ctx context.Context, id int) (core.Keepers, error) {
	return s.keeperStore.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, keeper core.Keepers) error {
	if keeper.CreatedAt.IsZero() {
		keeper.CreatedAt = time.Now()
	}

	return s.keeperStore.Create(ctx, keeper)
}

func (s *service) DeleteByID(ctx context.Context, id int) error {
	return s.keeperStore.DeleteByID(ctx, id)
}

func (s *service) SoftDeleteByID(ctx context.Context, id int) error {
	return s.keeperStore.SoftDeleteByID(ctx, id)
}

func (s *service) UpdateByID(ctx context.Context, keeper core.Keepers) (core.Keepers, error) {
	return s.keeperStore.UpdateByID(ctx, keeper)
}
