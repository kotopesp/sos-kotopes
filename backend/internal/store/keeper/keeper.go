package keepestore

import (
	"context"
	"fmt"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.KeeperStore {
	return &store{pg}
}

func (s *store) Create(ctx context.Context, keeper core.Keepers) error {
	if err := s.DB.WithContext(ctx).Create(&keeper).Error; err != nil {
		return err
	}
	return nil
}

func (s *store) DeleteByID(ctx context.Context, id int) error {
	result := s.DB.WithContext(ctx).Delete(core.Keepers{}, id)
	if result.RowsAffected == 0 {
		return core.ErrRecordNotFound
	}

	return result.Error
}

func (s *store) UpdateByID(ctx context.Context, keeper core.Keepers) (core.Keepers, error) {
	result := s.DB.WithContext(ctx).Model(&core.Keepers{}).Where("id = ?", keeper.ID).Updates(keeper)
	if result.RowsAffected == 0 {
		return keeper, core.ErrRecordNotFound
	}

	return keeper, result.Error
}

func (s *store) GetAll(ctx context.Context, params core.GetAllKeepersParams) ([]core.Keepers, error) {
	var keepers []core.Keepers
	query := s.DB.WithContext(ctx).Model(&core.Keepers{})

	if params.Location != nil {
		query = query.Where("location LIKE = ?", "%"+*params.Location+"%")
	}

	if params.MinPrice != nil {
		query = query.Where("price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("price <= ?", *params.MaxPrice)
	}

	query = query.Select("keepers.*, AVG(keeper_reviews.grade) as avg_grade").
		Joins("left join keeper_reviews on keeper_reviews.keeper_id = keepers.id").
		Group("keepers.id")

	if params.MinRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) >= ?", *params.MinRating)
	}
	if params.MaxRating != nil {
		query = query.Having("AVG(keeper_reviews.grade) <= ?", *params.MaxRating)
	}

	if params.SortBy != nil {
		sortField := *params.SortBy
		sortOrder := "asc"
		if params.SortOrder != nil && (*params.SortOrder == "desc" || *params.SortOrder == "DESC") {
			sortOrder = "desc"
		}

		if sortField == "avg_grade" || sortField == "price" || sortField == "created_at" {
			query = query.Order(sortField + " " + sortOrder)
		}
	}

	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	if err := query.Find(&keepers).Error; err != nil {
		return nil, fmt.Errorf("keepers.GetAll.find(): %s", err.Error())
	}

	return keepers, nil
}

func (s *store) GetByID(ctx context.Context, id int) (core.Keepers, error) {
	var keeper = core.Keepers{ID: id}

	if err := s.DB.WithContext(ctx).First(&keeper).Error; err != nil {
		return core.Keepers{}, err
	}

	return keeper, nil
}
