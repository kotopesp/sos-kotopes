package core

import (
	"context"
	"time"
)

type Keepers struct {
	ID                   int       `gorm:"primaryKey;autoIncrement;column:id"`
	UserID               int       `gorm:"column:user_id"`
	Description          string    `gorm:"column:description"`
	Price                float64   `gorm:"column:price"`
	Location             string    `gorm:"column:location"`
	HasCage              bool      `gorm:"column:has_cage"`
	BoardingDuration     string    `gorm:"column:boarding_duration"`
	BoardingCompensation string    `gorm:"column:boarding_compensation"`
	AnimalAcceptance     string    `gorm:"column:animal_acceptance"`
	AnimalCategory       string    `gorm:"column:animal_category"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime;column:updated_at"`
	IsDeleted            bool      `gorm:"column:is_deleted"`
	DeletedAt            time.Time `gorm:"column:deleted_at"`
}

type UpdateKeepers struct {
	ID                   int       `gorm:"primaryKey;autoIncrement;column:id"`
	UserID               int       `gorm:"column:user_id"`
	Description          string    `gorm:"column:description"`
	Price                float64   `gorm:"column:price"`
	Location             string    `gorm:"column:location"`
	HasCage              bool      `gorm:"column:has_cage"`
	BoardingDuration     string    `gorm:"column:boarding_duration"`
	BoardingCompensation string    `gorm:"column:boarding_compensation"`
	AnimalAcceptance     string    `gorm:"column:animal_acceptance"`
	AnimalCategory       string    `gorm:"column:animal_category"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type KeepersDetails struct {
	Keeper Keepers
	User   User
}

type GetAllKeepersParams struct {
	SortBy               *string
	SortOrder            *string
	Location             *string
	MinRating            *float64
	MaxRating            *float64
	MinPrice             *float64
	MaxPrice             *float64
	HasCage              *bool
	BoardingDuration     *string
	BoardingCompensation *string
	AnimalAcceptance     *string
	AnimalCategory       *string
	Limit                *int
	Offset               *int
}

type KeeperStore interface {
	GetAll(ctx context.Context, params GetAllKeepersParams) ([]Keepers, error)
	GetByID(ctx context.Context, id int) (Keepers, error)
	Create(ctx context.Context, keeper Keepers) error
	SoftDeleteByID(ctx context.Context, id int) error
	UpdateByID(ctx context.Context, keeper UpdateKeepers) (Keepers, error)
}

type KeeperService interface {
	GetAll(ctx context.Context, params GetAllKeepersParams) ([]KeepersDetails, error)
	GetByID(ctx context.Context, id int) (KeepersDetails, error)
	Create(ctx context.Context, keeper Keepers) error
	SoftDeleteByID(ctx context.Context, id int, userID int) error
	UpdateByID(ctx context.Context, keeper UpdateKeepers) (KeepersDetails, error)

	KeeperReviewsService
}

// TableName table name in db for gorm
func (Keepers) TableName() string {
	return "keepers"
}
