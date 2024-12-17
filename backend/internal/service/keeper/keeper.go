package keeperservice

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	keeperStore       core.KeeperStore
	keeperReviewStore core.KeeperReviewStore
	userStore         core.UserStore
}

func New(keeperStore core.KeeperStore, keeperReviewStore core.KeeperReviewStore, userStore core.UserStore) core.KeeperService {
	return &service{
		keeperStore:       keeperStore,
		keeperReviewStore: keeperReviewStore,
		userStore:         userStore,
	}
}

func (s *service) GetAllKeepers(ctx context.Context, params core.GetAllKeepersParams) (data []core.Keeper, err error) {
	if *params.SortBy == "" {
		*params.SortBy = "created_at"
	}
	if *params.SortOrder == "" {
		*params.SortOrder = "desc"
	}

	keepers, err := s.keeperStore.GetAllKeepers(ctx, params)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return keepers, nil
}

func (s *service) GetKeepeByID(ctx context.Context, id int) (data core.Keeper, err error) {
	keeper, err := s.keeperStore.GetKeeperByID(ctx, id)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.Keeper{}, err
	}

	return keeper, nil
}

func (s *service) CreateKeeper(ctx context.Context, keeper core.Keeper) (data core.Keeper, err error) {
	keeper.CreatedAt = time.Now()
	keeper.IsDeleted = false
	keeper.UpdatedAt = time.Now()

	return keeper, s.keeperStore.CreateKeeper(ctx, keeper)
}

func (s *service) DeleteKeeper(ctx context.Context, id, userID int) error {
	storedKeeper, err := s.keeperStore.GetKeeperByID(ctx, id)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	if storedKeeper.UserID != userID {
		logger.Log().Debug(ctx, core.ErrKeeperUserIDMissmatch.Error())
		return core.ErrKeeperUserIDMissmatch
	} else if storedKeeper.IsDeleted {
		logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
		return core.ErrRecordNotFound
	}

	return s.keeperStore.DeleteKeeper(ctx, id)
}

func (s *service) UpdateKeeper(ctx context.Context, id, userID int, keeper core.UpdateKeeper) (core.Keeper, error) {
	storedKeeper, err := s.keeperStore.GetKeeperByID(ctx, id)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.Keeper{}, err
	}

	if storedKeeper.UserID != userID {
		logger.Log().Debug(ctx, core.ErrKeeperUserIDMissmatch.Error())
		return core.Keeper{}, core.ErrKeeperUserIDMissmatch
	} else if storedKeeper.IsDeleted {
		logger.Log().Debug(ctx, core.ErrRecordNotFound.Error())
		return core.Keeper{}, core.ErrRecordNotFound
	}

	updatedKeeper, err := s.keeperStore.UpdateKeeper(ctx, id, keeper)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return core.Keeper{}, err
	}

	return updatedKeeper, nil
}
