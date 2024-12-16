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
	query := s.DB.WithContext(ctx).Model(&core.Vets{})

	// Filter by price range
	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	// Add optional filters for ratings
	//if params.MinRating != nil {
	//	query = query.Having("AVG(vet_reviews.grade) >= ?", *params.MinRating)
	//}
	//if params.MaxRating != nil {
	//	query = query.Having("AVG(vet_reviews.grade) <= ?", *params.MaxRating)
	//}

	// Add optional location filter
	if params.Location != nil {
		query = query.Where("location = ?", *params.Location)
	}

	// Select and join to calculate the average grade
	//query = query.Select("vets.*, AVG(vet_reviews.grade) as avg_grade"). // Изменено на vets
	//	Joins("left join vet_reviews on vet_reviews.vet_id = vets.id").  // Изменено на vets
	//	Group("vets.id") // Изменено на vets

	//if params.SortBy != nil && params.SortOrder != nil {
	//	query = query.Order(*params.SortBy + " " + *params.SortOrder)
	//}

	//if params.Limit != nil {
	//	query = query.Limit(*params.Limit)
	//}
	//if params.Offset != nil {
	//	query = query.Offset(*params.Offset)
	//}

	if err := query.Find(&vets).Error; err != nil {
		logger.Log().Debug(ctx, err.Error())
		return nil, err
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

func (s *store) GetByOrgName(ctx context.Context, orgName string) (core.Vets, error) {
	var vet core.Vets

	if err := s.DB.WithContext(ctx).Where("org_name = ?", orgName).First(&vet).Error; err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			logger.Log().Error(ctx, core.ErrRecordNotFound.Error())
			return core.Vets{}, core.ErrRecordNotFound
		}

		logger.Log().Error(ctx, err.Error())
		return core.Vets{}, err
	}

	return vet, nil
}
