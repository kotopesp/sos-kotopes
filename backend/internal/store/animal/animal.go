package animal

import (
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"context"
	"errors"
	"time"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.AnimalStore {
	return &store{pg}
}

func (s *store) CreateAnimal(ctx context.Context, animal core.Animal) (core.Animal, error) {
	animal.CreatedAt = time.Now()
	animal.UpdatedAt = time.Now()
	if err := s.DB.WithContext(ctx).Create(&animal).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
        return core.Animal{}, err
    }

    return animal, nil
}

func (s *store) GetAnimalByID(ctx context.Context, id int) (core.Animal, error) {
	var animal core.Animal

	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&animal).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Animal{}, core.ErrPostNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Animal{}, err
	}

	return animal, nil
}

func (s *store) UpdateAnimal(ctx context.Context, animal core.Animal) error {
	animal.UpdatedAt = time.Now()
	if err := s.DB.WithContext(ctx).Save(&animal).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}
	
	return nil
}
