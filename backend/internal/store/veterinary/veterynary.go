package veterinarystore

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

func New(pg *postgres.Postgres) core.VeterinaryStore {
	return &store{pg}
}

func (s *store) Create(ctx context.Context, veterinary core.Veterinary) error {
	if err := s.DB.WithContext(ctx).Create(&veterinary).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteByID(ctx context.Context, id int) error {
	err := s.DB.WithContext(ctx).Model(&core.Veterinary{}).Where("id = ?", id).Updates(core.Veterinary{IsDeleted: true, DeletedAt: time.Now()}).Error

	return err
}

func (s *store) UpdateByID(ctx context.Context, veterinary core.UpdateVeterinary) (core.Veterinary, error) {
	veterinary.UpdatedAt = time.Now()

	var updatedVeterinary core.Veterinary

	result := s.DB.WithContext(ctx).Model(&core.Veterinary{}).Where("id = ? AND is_deleted = ?", veterinary.ID, false).Updates(veterinary).First(&updatedVeterinary, veterinary.ID)
	if result.Error != nil {

		if errors.Is(result.Error, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Veterinary{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, result.Error.Error())
		return core.Veterinary{}, result.Error
	}

	return updatedVeterinary, nil
}

func (s *store) GetAll(ctx context.Context, params core.GetAllVeterinaryParams) ([]core.Veterinary, error) {
	var veterinaries []core.Veterinary
	query := s.DB.WithContext(ctx).Model(&core.Veterinary{})

	// Filter by price range
	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	// Add optional filters for ratings
	if params.MinRating != nil {
		query = query.Having("AVG(veterinary_reviews.grade) >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Having("AVG(veterinary_reviews.grade) <= ?", *params.MaxRating)
	}

	// Add optional location filter
	if params.Location != nil {
		query = query.Where("location = ?", *params.Location)
	}

	// Select and join to calculate the average grade
	query = query.Select("veterinaries.*, AVG(veterinary_reviews.grade) as avg_grade").
		Joins("left join veterinary_reviews on veterinary_reviews.veterinary_id = veterinaries.id").
		Group("veterinaries.id")

	if params.SortBy != nil && params.SortOrder != nil {
		query = query.Order(*params.SortBy + " " + *params.SortOrder)
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&veterinaries).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
	}

	return veterinaries, nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Veterinary, error) {
	var veterinary = core.Veterinary{ID: id}

	if err := s.DB.WithContext(ctx).First(&veterinary).Error; err != nil {
		return core.Veterinary{}, err
	}

	return veterinary, nil
}

func (s *store) GetByOrgName(ctx context.Context, orgName string) (core.Veterinary, error) {
	var veterinary core.Veterinary

	if err := s.DB.WithContext(ctx).Where("org_name = ?", orgName).First(&veterinary).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Veterinary{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Veterinary{}, err
	}

	return veterinary, nil
}
