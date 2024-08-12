package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	keeperStore        core.KeeperStore
	keeperReviewsStore core.KeeperReviewsStore
	userStore          core.UserStore
}

func New(keeperStore core.KeeperStore, keeperReviewStore core.KeeperReviewsStore, userStore core.UserStore) core.KeeperService {
	return &service{
		keeperStore:        keeperStore,
		keeperReviewsStore: keeperReviewStore,
		userStore:          userStore,
	}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.KeepersDetails, error) {
	if *params.SortBy == "" {
		*params.SortBy = "created_at"
	}
	if *params.SortOrder == "" {
		*params.SortBy = "desc"
	}

	keepers, err := s.keeperStore.GetAll(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	keepersDetails := make([]core.KeepersDetails, len(keepers))

	for i, keeper := range keepers {
		keeperUser, err := s.userStore.GetUserByID(ctx, keeper.UserID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		keepersDetails[i] = core.KeepersDetails{
			Keeper: keeper,
			User:   keeperUser,
		}
	}

	return keepersDetails, nil
}

func (s *service) GetByID(ctx context.Context, id int) (core.KeepersDetails, error) {
	keeper, err := s.keeperStore.GetByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeepersDetails{}, err
	}

	keeperUser, err := s.userStore.GetUserByID(ctx, keeper.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeepersDetails{}, err
	}

	return core.KeepersDetails{
		Keeper: keeper,
		User:   keeperUser,
	}, nil
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

func (s *service) UpdateByID(ctx context.Context, keeper core.UpdateKeepers) (core.KeepersDetails, error) {
	updatedKeeper, err := s.keeperStore.UpdateByID(ctx, keeper)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeepersDetails{}, err
	}

	keeperUser, err := s.userStore.GetUserByID(ctx, updatedKeeper.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.KeepersDetails{}, err
	}

	return core.KeepersDetails{
		Keeper: updatedKeeper,
		User:   keeperUser,
	}, nil
}
