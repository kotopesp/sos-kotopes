package veterinary

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"time"
)

type service struct {
	VeterinaryStore core.VeterinaryStore
	userStore       core.UserStore
}

func New(VeterinaryStore core.VeterinaryStore, userStore core.UserStore) core.VeterinaryService {
	return &service{
		VeterinaryStore: VeterinaryStore,
		userStore:       userStore,
	}
}

func (s *service) GetAll(ctx context.Context, params core.GetAllVeterinaryParams) ([]core.VeterinaryDetails, error) {
	if *params.SortBy == "" {
		*params.SortBy = "created_at"
	}
	if *params.SortOrder == "" {
		*params.SortOrder = "desc"
	}

	veterinary, err := s.VeterinaryStore.GetAll(ctx, params)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	VeterinaryDetails := make([]core.VeterinaryDetails, len(veterinary))

	for i, veterinary := range veterinary {
		veterinaryUser, err := s.userStore.GetUserByID(ctx, veterinary.UserID)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		VeterinaryDetails[i] = core.VeterinaryDetails{
			Veterinary: veterinary,
			User:       veterinaryUser,
		}
	}

	return VeterinaryDetails, nil
}

func (s *service) GetByID(ctx context.Context, id int) (core.VeterinaryDetails, error) {
	veterinary, err := s.VeterinaryStore.GetByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	veterinaryUser, err := s.userStore.GetUserByID(ctx, veterinary.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	return core.VeterinaryDetails{
		Veterinary: veterinary,
		User:       veterinaryUser,
	}, nil
}

// переделать под то что нада
func (s *service) GetByOrgName(ctx context.Context, id int) (core.VeterinaryDetails, error) {
	veterinary, err := s.VeterinaryStore.GetByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	veterinaryUser, err := s.userStore.GetUserByID(ctx, veterinary.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	return core.VeterinaryDetails{
		Veterinary: veterinary,
		User:       veterinaryUser,
	}, nil
}

func (s *service) Create(ctx context.Context, veterinary core.Veterinary) error {
	if veterinary.CreatedAt.IsZero() {
		veterinary.CreatedAt = time.Now()
	}

	return s.VeterinaryStore.Create(ctx, veterinary)
}

func (s *service) DeleteByID(ctx context.Context, id, userID int) error {
	storedVeterinary, err := s.VeterinaryStore.GetByID(ctx, id)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	if storedVeterinary.UserID != userID {
		logger.Log().Error(ctx, core.ErrKeeperUserIDMissmatch.Error())
		return core.ErrKeeperUserIDMissmatch
	} else if storedVeterinary.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.ErrRecordNotFound
	}

	return s.VeterinaryStore.DeleteByID(ctx, id)
}

func (s *service) UpdateByID(ctx context.Context, veterinary core.UpdateVeterinary) (core.VeterinaryDetails, error) {
	storedVeterinary, err := s.VeterinaryStore.GetByID(ctx, veterinary.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	if storedVeterinary.UserID != veterinary.UserID {
		logger.Log().Error(ctx, core.ErrKeeperUserIDMissmatch.Error())
		return core.VeterinaryDetails{}, core.ErrKeeperUserIDMissmatch
	} else if storedVeterinary.IsDeleted {
		logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
		return core.VeterinaryDetails{}, core.ErrRecordNotFound
	}

	updatedVeterinary, err := s.VeterinaryStore.UpdateByID(ctx, veterinary) //тут ещё раз проверить
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	veterinaryUser, err := s.userStore.GetUserByID(ctx, updatedVeterinary.UserID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.VeterinaryDetails{}, err
	}

	return core.VeterinaryDetails{
		Veterinary: updatedVeterinary,
		User:       veterinaryUser,
	}, nil
}
