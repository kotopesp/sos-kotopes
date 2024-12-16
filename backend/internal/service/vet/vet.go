package vet

import (
	"context"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	vetStore        core.VetStore
	vetReviewsStore core.VetReviewsStore
	userStore       core.UserStore
}

func New(vetStore core.VetStore, vetReviewsStore core.VetReviewsStore, userStore core.UserStore) core.VetService {
	return &service{
		vetStore:        vetStore,
		vetReviewsStore: vetReviewsStore,
		userStore:       userStore,
	}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllVetParams) ([]core.VetsDetails, error) {
	vet, err := s.vetStore.GetAll(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	var vetsDetails []core.VetsDetails

	for _, v := range vet {
		if v.IsDeleted {
			continue
		}

		vetUser, err := s.userStore.GetUserByID(ctx, v.UserID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		vetsDetails = append(vetsDetails, core.VetsDetails{
			Vet:  v,
			User: vetUser,
		})
	}

	return vetsDetails, nil
}

func (s *service) GetByUserID(ctx context.Context, userID int) (core.VetsDetails, error) {
	vet, err := s.vetStore.GetByUserID(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetsDetails{}, err
	}

	vetUser, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetsDetails{}, err
	}

	return core.VetsDetails{
		Vet:  vet,
		User: vetUser,
	}, nil
}

func (s *service) Create(ctx context.Context, vet core.Vets) error {
	if vet.CreatedAt.IsZero() {
		vet.CreatedAt = time.Now()
	}

	return s.vetStore.Create(ctx, vet)
}

func (s *service) DeleteByUserID(ctx context.Context, userID int) error {
	storedVet, err := s.vetStore.GetByUserID(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	if storedVet.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.ErrRecordNotFound
	}

	return s.vetStore.DeleteByUserID(ctx, storedVet.UserID)
}

func (s *service) UpdateByUserID(ctx context.Context, vet core.UpdateVets) (core.VetsDetails, error) {
	storedVet, err := s.vetStore.GetByUserID(ctx, vet.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetsDetails{}, err
	}

	if storedVet.UserID != vet.UserID {
		logger.Log().Error(ctx, core.ErrVetUserIDMismatch.Error())
		return core.VetsDetails{}, core.ErrVetUserIDMismatch
	} else if storedVet.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.VetsDetails{}, core.ErrRecordNotFound
	}

	updatedVet, err := s.vetStore.UpdateByUserID(ctx, vet)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetsDetails{}, err
	}

	vetUser, err := s.userStore.GetUserByID(ctx, updatedVet.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VetsDetails{}, err
	}

	return core.VetsDetails{
		Vet:  updatedVet,
		User: vetUser,
	}, nil
}
