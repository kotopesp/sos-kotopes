package core

import (
	"context"
	"time"
)

type Keepers struct {
	ID          int       `gorm:"primaryKey;autoIncrement;column:id"`
	UserID      int       `gorm:"column:user_id"`
	Description string    `gorm:"column:description"`
	Price       float64   `gorm:"column:price"`
	Location    string    `gorm:"column:location"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type GetAllKeepersParams struct {
	SortBy    *string
	SortOrder *string
	Location  *string
	MinRating *float64
	MaxRating *float64
	MinPrice  *float64
	MaxPrice  *float64
	Limit     *int
	Offset    *int
}

type KeeperStore interface {
	GetAll(ctx context.Context, params GetAllKeepersParams) ([]Keepers, error)
	GetByID(ctx context.Context, id int) (Keepers, error)
	Create(ctx context.Context, keeper Keepers) error
	DeleteByID(ctx context.Context, id int) error
	UpdateByID(ctx context.Context, keeper Keepers) error
}

type KeeperService interface {
	GetAll(ctx context.Context, params GetAllKeepersParams) ([]Keepers, error)
	GetByID(ctx context.Context, id int) (Keepers, error)
	Create(ctx context.Context, keeper Keepers) error
	DeleteByID(ctx context.Context, id int) error
	UpdateByID(ctx context.Context, keeper Keepers) error

	KeeperReviewsService
}

// TableName table name in db for gorm
func (Keepers) TableName() string {
	return "keepers"
}
