package vetstore

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"time"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.VetStore {
	return &store{pg}
}

func (s *store) Create(ctx context.Context, vet core.Vets) error {
	if err := s.DB.WithContext(ctx).Create(&vet).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteByUserID(ctx context.Context, userID int) error {
	err := s.DB.WithContext(ctx).
		Model(&core.Vets{}).
		Where("user_id = ?", userID).
		Updates(core.Vets{
			IsDeleted: true,
			DeletedAt: time.Now(),
		}).Error

	return err
}

func (s *store) UpdateByUserID(ctx context.Context, vet core.UpdateVets) (core.Vets, error) {
	vet.UpdatedAt = time.Now()

	var updatedVet core.Vets

	result := s.DB.WithContext(ctx).Model(&core.Vets{}).Where("id = ? AND is_deleted = ?", vet.UserID, false).Updates(vet).First(&updatedVet, vet.UserID)
	if result.Error != nil {

		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Vets{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, result.Error.Error())
		return core.Vets{}, result.Error
	}

	return updatedVet, nil
}

func (s *store) GetAll(ctx context.Context, params core.GetAllVetParams) ([]core.Vets, error) {
	var vets []core.Vets

	query := s.DB.WithContext(ctx).Model(&core.Vets{}).Where("vets.is_deleted = ?", false)

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	logger.Log().Debug(ctx, "params", params)

	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	logger.Log().Debug(ctx, "after price query:", query.Statement.SQL.String())

	if params.MinRating != nil {
		query = query.Having("AVG(vet_reviews.grade) >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Having("AVG(vet_reviews.grade) <= ?", *params.MaxRating)
	}

	if params.Location != nil {
		query = query.Where("location = ?", *params.Location)
	}

	query = query.Select("vets.*, AVG(vet_reviews.grade) as avg_grade").
		Joins("LEFT JOIN vet_reviews ON vet_reviews.vet_id = vets.id").
		Group("vets.id")

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&vets).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return []core.Vets{}, nil
	}
	return vets, nil
}

func (s *store) GetByUserID(ctx context.Context, userID int) (core.Vets, error) {
	var vet core.Vets

	if err := s.DB.WithContext(ctx).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		First(&vet).Error; err != nil {
		return core.Vets{}, err
	}
	return vet, nil
}
