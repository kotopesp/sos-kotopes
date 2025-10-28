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
		LocationID       int       `gorm:"column:location_id"`
		EquipmentRental  int       `gorm:"column:equipment_rental"`
		HaveMetalCage    bool      `gorm:"column:have_metal_cage"`
		HavePlasticCage  bool      `gorm:"column:have_plastic_cage"`
		HaveNet          bool      `gorm:"column:have_net"`
		HaveLadder       bool      `gorm:"column:have_ladder"`
		HaveOther        string    `gorm:"column:have_other"`
		Price            int       `gorm:"column:price"`
		HaveCar          bool      `gorm:"column:have_car"`
		WillingnessCarry string    `gorm:"column:willingness_carry"`
		CreatedAt        time.Time `gorm:"column:created_at"`
		UpdatedAt        time.Time `gorm:"column:updated_at"`
		IsDeleted        bool      `gorm:"column:is_deleted"`
		DeletedAt        time.Time `gorm:"column:deleted_at"`
	}

	UpdateSeeker struct {
		ID               *int      `gorm:"column:id"`
		UserID           *int      `gorm:"column:user_id"`
		AnimalType       *string   `gorm:"column:animal_type"`
		Description      *string   `gorm:"column:description"`
		LocationID       *int      `gorm:"column:location_id"`
		EquipmentRental  *int      `gorm:"column:equipment_rental"`
		HaveMetalCage    *bool     `gorm:"column:have_metal_cage"`
		HavePlasticCage  *bool     `gorm:"column:have_plastic_cage"`
		HaveNet          *bool     `gorm:"column:have_net"`
		HaveLadder       *bool     `gorm:"column:have_ladder"`
		HaveOther        *string   `gorm:"column:have_other"`
		Price            *int      `gorm:"column:price"`
		HaveCar          *bool     `gorm:"column:have_car"`
		WillingnessCarry *string   `gorm:"column:willingness_carry"`
		UpdatedAt        time.Time `gorm:"autoUpdateTime;column:updated_at"`
	}

	DeleteSeeker struct {
		UserID int `gorm:"column:user_id"`
	}

	GetAllSeekersParams struct {
		SortBy             *string
		SortOrder          *string
		AnimalType         *string
		LocationID         *int
		MinEquipmentRental *int
		MaxEquipmentRental *int
		HaveMetalCage      *bool
		HavePlasticCage    *bool
		HaveNet            *bool
		HaveLadder         *bool
		HaveOther          *string
		MinPrice           *int
		MaxPrice           *int
		HaveCar            *bool
		Limit              *int
		Offset             *int
	}

	SeekersService interface {
		CreateSeeker(ctx context.Context, seeker Seeker) (Seeker, error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeeker) (Seeker, error)
		DeleteSeeker(ctx context.Context, userID int) error
		GetAllSeekers(ctx context.Context, params GetAllSeekersParams) ([]Seeker, error)
	}

	SeekersStore interface {
		CreateSeeker(ctx context.Context, seeker Seeker) (createSeeker Seeker, err error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, userID int, updateSeeker map[string]interface{}) (seeker Seeker, err error)
		DeleteSeeker(ctx context.Context, userID int) error
		GetAllSeekers(ctx context.Context, params GetAllSeekersParams) ([]Seeker, error)
	}
)
