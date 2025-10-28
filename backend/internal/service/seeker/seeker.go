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

func (s *service) CreateSeeker(ctx context.Context, seeker core.Seeker) (createSeeker core.Seeker, err error) {
	if _, err = s.seekersStore.GetSeeker(ctx, seeker.UserID); !errors.Is(err, core.ErrSeekerNotFound) {
		logger.Log().Error(ctx, core.ErrSeekerExists.Error())
		return core.Seeker{}, core.ErrSeekerExists
	}

	return s.seekersStore.CreateSeeker(ctx, seeker)
}

func (s *service) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	seeker, err := s.seekersStore.GetSeeker(ctx, id)
	if err != nil || seeker.IsDeleted {
		logger.Log().Error(ctx, core.ErrSeekerNotFound.Error())
		return core.Seeker{}, core.ErrSeekerNotFound
	}

	return seeker, nil
}

func createMapUpdates(updateSeeker core.UpdateSeeker, seeker core.Seeker) map[string]interface{} {
	updates := make(map[string]interface{})

	if updateSeeker.AnimalType != nil && seeker.AnimalType != *updateSeeker.AnimalType {
		updates["animal_type"] = *updateSeeker.AnimalType
	}

	if updateSeeker.Description != nil && seeker.Description != *updateSeeker.Description {
		updates["description"] = *updateSeeker.Description
	}

	if updateSeeker.LocationId != nil && seeker.LocationId != *updateSeeker.LocationId {
		updates["location_id"] = *updateSeeker.LocationId
	}

	if updateSeeker.EquipmentRental != nil && seeker.EquipmentRental != *updateSeeker.EquipmentRental {
		updates["equipment_rental"] = *updateSeeker.EquipmentRental
	}

	if updateSeeker.HaveMetalCage != nil && seeker.HaveMetalCage != *updateSeeker.HaveMetalCage {
		updates["have_metal_cage"] = *updateSeeker.HaveMetalCage
	}

	if updateSeeker.HavePlasticCage != nil && seeker.HavePlasticCage != *updateSeeker.HavePlasticCage {
		updates["have_plastic_cage"] = *updateSeeker.HavePlasticCage
	}

	if updateSeeker.HaveNet != nil && seeker.HaveNet != *updateSeeker.HaveNet {
		updates["have_net"] = *updateSeeker.HaveNet
	}

	if updateSeeker.HaveLadder != nil && seeker.HaveLadder != *updateSeeker.HaveLadder {
		updates["have_ladder"] = *updateSeeker.HaveLadder
	}

	if updateSeeker.HaveOther != nil && seeker.HaveOther != *updateSeeker.HaveOther {
		updates["have_other"] = *updateSeeker.HaveOther
	}

	if updateSeeker.Price != nil && seeker.Price != *updateSeeker.Price {
		updates["price"] = *updateSeeker.Price
	}

	if updateSeeker.HaveCar != nil && seeker.HaveCar != *updateSeeker.HaveCar {
		updates["have_car"] = *updateSeeker.HaveCar
	}

	if updateSeeker.WillingnessCarry != nil && seeker.WillingnessCarry != *updateSeeker.WillingnessCarry {
		updates["willingness_carry"] = *updateSeeker.WillingnessCarry
	}

	return updates
}

func (s *service) UpdateSeeker(ctx context.Context, updateSeeker core.UpdateSeeker) (core.Seeker, error) {
	getSeeker, err := s.GetSeeker(ctx, *updateSeeker.UserID)
	if err != nil {
		logger.Log().Error(ctx, core.ErrSeekerNotFound.Error())
		return core.Seeker{}, core.ErrSeekerNotFound
	}

	mapUpdates := createMapUpdates(updateSeeker, getSeeker)
	if len(mapUpdates) == 0 {
		logger.Log().Error(ctx, core.ErrEmptyUpdateRequest.Error())
		return core.Seeker{}, core.ErrEmptyUpdateRequest
	}

	seeker, err := s.seekersStore.UpdateSeeker(ctx, getSeeker.UserID, mapUpdates)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Seeker{}, err
	}

	return seeker, nil
}

func (s *service) DeleteSeeker(ctx context.Context, userID int) error {
	seeker, err := s.GetSeeker(ctx, userID)
	if err != nil || seeker.IsDeleted {
		logger.Log().Error(ctx, core.ErrSeekerNotFound.Error())
		return core.ErrSeekerNotFound
	}

	if err = s.seekersStore.DeleteSeeker(ctx, userID); err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func validateParams(params core.GetAllSeekersParams) core.GetAllSeekersParams {
	sortBy, sortOrder := "created_at", "desc"
	if params.SortBy == nil {
		params.SortBy = &sortBy
	}
	if params.SortOrder == nil {
		params.SortOrder = &sortOrder
	}

	limit, offset := 10, 0
	if params.Limit == nil {
		params.Limit = &limit
	}
	if params.Offset == nil {
		params.Offset = &offset
	}

	return params
}

func (s *service) GetAllSeekers(ctx context.Context, params core.GetAllSeekersParams) ([]core.Seeker, error) {
	params = validateParams(params)
	seekers, err := s.seekersStore.GetAllSeekers(ctx, params)
	if err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return seekers, nil
}
