package animal

import (
	"context"
	"errors"
	"time"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type (
	store struct {
		*postgres.Postgres
	}
)

func New(pg *postgres.Postgres) core.AnimalStore {
	return &store{pg}
}

// CreateAnimal inserts a new animal record into the database
func (s *store) CreateAnimal(ctx context.Context, animal core.Animal) (core.Animal, error) {
	// Set the creation and update timestamps
	animal.CreatedAt = time.Now().UTC()
	animal.UpdatedAt = time.Now().UTC()

	if err := s.DB.WithContext(ctx).Create(&animal).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return core.Animal{}, err
	}

	return animal, nil
}

// GetAnimalByID retrieves an animal record from the database by its ID
func (s *store) GetAnimalByID(ctx context.Context, id int) (core.Animal, error) {
	var animal core.Animal

	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&animal).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Animal{}, core.ErrAnimalNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Animal{}, err
	}

	return animal, nil
}

// UpdateAnimal updates an existing animal record in the database
func (s *store) UpdateAnimal(ctx context.Context, animal core.Animal) (core.Animal, error) {
	// Set the update timestamp
	animal.UpdatedAt = time.Now().UTC()

	var updateAnimal core.Animal

	if err := s.DB.WithContext(ctx).Save(&animal).First(&updateAnimal).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Animal{}, core.ErrAnimalNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Animal{}, err
	}

	return updateAnimal, nil
}
