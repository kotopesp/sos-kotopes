package core

import (
	"context"
	"time"
)

type Keeper struct {
	ID                   int        `gorm:"primaryKey;autoIncrement;column:id"`
	UserID               int        `gorm:"column:user_id"`
	User                 User       `gorm:"foreignKey:UserID;references:ID"`
	Description          *string    `gorm:"column:description"`
	Price                *float64   `gorm:"column:price"`
	LocationID           *int       `gorm:"column:location_id"`
	HasCage              bool       `gorm:"column:has_cage"`
	BoardingDuration     string     `gorm:"column:boarding_duration"`
	BoardingCompensation string     `gorm:"column:boarding_compensation"`
	AnimalAcceptance     string     `gorm:"column:animal_acceptance"`
	AnimalCategory       string     `gorm:"column:animal_category"`
	CreatedAt            time.Time  `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime;column:updated_at"`
	IsDeleted            bool       `gorm:"column:is_deleted"`
	DeletedAt            *time.Time `gorm:"column:deleted_at"`
}

type GetAllKeepersParams struct {
	SortBy               *string
	SortOrder            *string
	LocationID           *int
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
	GetAllKeepers(ctx context.Context, params GetAllKeepersParams) (data []Keeper, err error)
	GetKeeperByID(ctx context.Context, id int) (Keeper, error)
	CreateKeeper(ctx context.Context, keeper Keeper) (data Keeper, err error)
	DeleteKeeper(ctx context.Context, id int) error
	UpdateKeeper(ctx context.Context, id int, keeper Keeper) (Keeper, error)
}

type KeeperService interface {
	GetAllKeepers(ctx context.Context, params GetAllKeepersParams) (data []Keeper, err error)
	GetKeepeByID(ctx context.Context, id int) (data Keeper, err error)
	CreateKeeper(ctx context.Context, keeper Keeper) (data Keeper, err error)
	DeleteKeeper(ctx context.Context, id int, userID int) error
	UpdateKeeper(ctx context.Context, id int, userID int, keeper Keeper) (Keeper, error)

	KeeperReviewService
}

// TableName table name in db for gorm
func (Keeper) TableName() string {
	return "keepers"
}
