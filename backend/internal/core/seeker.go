package core

import (
	"context"
	"time"
)

type (
	Seeker struct {
		ID               int       `gorm:"primaryKey;autoIncrement;column:id"`
		UserID           int       `gorm:"column:user_id"`
		User             User      `gorm:"foreignKey:UserID;references:ID"`
		AnimalType       string    `gorm:"column:animal_type"`
		Description      string    `gorm:"column:description"`
		Location         string    `gorm:"column:location"`
		EquipmentRental  int       `gorm:"column:equipment_rental"`
		EquipmentID      int       `gorm:"column:equipment_id"`
		Price            int       `gorm:"column:price"`
		HaveCar          bool      `gorm:"column:have_car"`
		WillingnessCarry string    `gorm:"column:willingness_carry"`
		CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
		UpdatedAt        time.Time `gorm:"autoUpdateTime;column:updated_at"`
		IsDeleted        bool      `gorm:"column:is_deleted"`
		DeletedAt        time.Time `gorm:"column:deleted_at"`
	}

	UpdateSeeker struct {
		ID               *int      `gorm:"column:id"`
		UserID           *int      `gorm:"column:user_id"`
		AnimalType       *string   `gorm:"column:animal_type"`
		Description      *string   `gorm:"column:description"`
		Location         *string   `gorm:"column:location"`
		EquipmentRental  *int      `gorm:"column:equipment_rental"`
		EquipmentID      *int      `gorm:"column:equipment_id"`
		Price            *int      `gorm:"column:price"`
		HaveCar          *bool     `gorm:"column:have_car"`
		WillingnessCarry *string   `gorm:"column:willingness_carry"`
		UpdatedAt        time.Time `gorm:"autoUpdateTime;column:updated_at"`
	}

	Equipment struct {
		ID              int    `gorm:"primaryKey;autoIncrement;column:id"`
		HaveMetalCage   bool   `gorm:"column:have_metal_cage"`
		HavePlasticCage bool   `gorm:"column:have_plastic_cage"`
		HaveNet         bool   `gorm:"column:have_net"`
		HaveLadder      bool   `gorm:"column:have_ladder"`
		HaveOther       string `gorm:"column:have_other"`
	}

	DeleteSeeker struct {
		UserID int `gorm:"column:user_id"`
	}

	GetAllSeekersParams struct {
		SortBy     *string
		SortOrder  *string
		AnimalType *string
		Location   *string
		Price      *int
		HaveCar    *bool
		Limit      *int
		Offset     *int
	}

	SeekersService interface {
		CreateSeeker(ctx context.Context, seeker Seeker, equipment Equipment) (Seeker, error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeeker) (Seeker, error)
		DeleteSeeker(ctx context.Context, userID int) error
		GetAllSeekers(ctx context.Context, params GetAllSeekersParams) ([]Seeker, error)
	}

	SeekersStore interface {
		CreateSeeker(ctx context.Context, seeker Seeker, equipment Equipment) (Seeker, error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, userID int, updateSeeker map[string]interface{}) (Seeker, error)
		DeleteSeeker(ctx context.Context, userID int) error
		GetAllSeekers(ctx context.Context, params GetAllSeekersParams) ([]Seeker, error)
	}
)

// TableName table name in db for gorm
func (Seeker) TableName() string {
	return "seekers"
}
