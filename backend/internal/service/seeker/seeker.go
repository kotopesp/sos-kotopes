package seeker

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	seekersStore core.SeekersStore
}

func New(seekersStore core.SeekersStore) core.SeekersService {
	return &service{seekersStore: seekersStore}
}

func (s *service) CreateSeeker(ctx context.Context, seeker core.Seeker) (core.Seeker, error) {
	if _, err := s.GetSeeker(ctx, seeker.UserID); !errors.Is(err, core.ErrSeekerNotFound) {
		logger.Log().Error(ctx, core.ErrSeekerExists.Error())
		return core.Seeker{}, core.ErrSeekerExists
	}

	return s.seekersStore.CreateSeeker(ctx, seeker)
}

func (s *service) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	seeker, err := s.seekersStore.GetSeeker(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, core.ErrSeekerNotFound.Error())
		return core.Seeker{}, core.ErrSeekerNotFound
	}

	if seeker.IsDeleted {
		logger.Log().Error(ctx, core.ErrSeekerDeleted.Error())
		return core.Seeker{}, core.ErrSeekerDeleted
	}

	return seeker, nil
}

func (s *service) UpdateSeeker(ctx context.Context, updateSeeker core.UpdateSeeker) (core.Seeker, error) {
	updates := make(map[string]interface{})

	if updateSeeker.AnimalType != nil {
		updates["animal_type"] = *updateSeeker.AnimalType
	}

	if updateSeeker.Description != nil {
		updates["description"] = *updateSeeker.Description
	}

	if updateSeeker.Location != nil {
		updates["location"] = *updateSeeker.Location
	}

	if updateSeeker.EquipmentRental != nil {
		updates["equipment_rental"] = *updateSeeker.EquipmentRental
	}

	if updateSeeker.HaveMetalCage != nil {
		updates["have_metal_cage"] = *updateSeeker.HaveMetalCage
	}

	if updateSeeker.HavePlasticCage != nil {
		updates["have_plastic_cage"] = *updateSeeker.HavePlasticCage
	}

	if updateSeeker.HaveNet != nil {
		updates["have_net"] = *updateSeeker.HaveNet
	}

	if updateSeeker.HaveLadder != nil {
		updates["have_ladder"] = *updateSeeker.HaveLadder
	}

	if updateSeeker.HaveOther != nil {
		updates["have_other"] = *updateSeeker.HaveOther
	}

	if updateSeeker.Price != nil {
		updates["price"] = *updateSeeker.Price
	}

	if updateSeeker.HaveCar != nil {
		updates["have_car"] = *updateSeeker.HaveCar
	}

	if updateSeeker.WillingnessCarry != nil {
		updates["willingness_carry"] = *updateSeeker.WillingnessCarry
	}

	if len(updates) == 0 {
		logger.Log().Error(ctx, core.ErrEmptyUpdateRequest.Error())
		return core.Seeker{}, core.ErrEmptyUpdateRequest
	}

	getSeeker, err := s.GetSeeker(ctx, *updateSeeker.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return s.seekersStore.UpdateSeeker(ctx, getSeeker.ID, updates)
}

func (s *service) DeleteSeeker(ctx context.Context, userID int) error {
	if _, err := s.GetSeeker(ctx, userID); err != nil {
		logger.Log().Debug(ctx, err.Error())
		return err
	}

	return s.seekersStore.DeleteSeeker(ctx, userID)
}

func (s *service) GetAllSeekers(ctx context.Context, params core.GetAllSeekersParams) ([]core.Seeker, error) {
	if params.SortBy == nil {
		params.SortBy = new(string)
	}
	if params.SortOrder == nil {
		params.SortOrder = new(string)
	}

	if *params.SortBy == "" {
		*params.SortBy = "created_at"
	}

	if *params.SortOrder == "" {
		*params.SortOrder = "desc"
	}

	seekers, err := s.seekersStore.GetAllSeekers(ctx, params)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return seekers, nil
}
