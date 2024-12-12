package core

import (
	"context"
	"time"
)

type (
	Seeker struct {
		ID          int       `gorm:"primaryKey;autoIncrement;column:id"`
		UserID      int       `gorm:"column:user_id"`
		Description string    `gorm:"column:description"`
		Location    string    `gorm:"column:location"`
		EquipmentID int       `gorm:"column:equipment_id"`
		Price       int       `gorm:"column:price"`
		HaveCar     bool      `gorm:"column:have_car"`
		CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
		IsDeleted   bool      `gorm:"column:is_deleted"`
		DeletedAt   time.Time `gorm:"column:deleted_at"`
	}

	UpdateSeeker struct {
		ID          *int      `gorm:"primaryKey;autoIncrement;column:id"`
		UserID      *int      `gorm:"column:user_id"`
		Description *string   `gorm:"column:description"`
		Location    *string   `gorm:"column:location"`
		EquipmentID *int      `gorm:"column:equipment_id"`
		Price       *int      `gorm:"column:price"`
		HaveCar     *bool     `gorm:"column:have_car"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
	}

	Equipment struct {
		ID              int    `gorm:"primaryKey;autoIncrement;column:id"`
		HaveMetalCage   bool   `gorm:"column:have_metal_cage"`
		HavePlasticCage bool   `gorm:"column:have_plastic_cage"`
		HaveNet         bool   `gorm:"column:have_net"`
		HaveLadder      bool   `gorm:"column:have_ladder"`
		HaveOther       string `gorm:"column:have_other"`
	}

	SeekersService interface {
		CreateSeeker(ctx context.Context, seeker Seeker, equipment Equipment) (Seeker, error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeeker) (Seeker, error)
	}

	SeekersStore interface {
		CreateSeeker(ctx context.Context, seeker Seeker, equipment Equipment) (Seeker, error)
		GetSeeker(ctx context.Context, userID int) (Seeker, error)
		UpdateSeeker(ctx context.Context, userID int, updateSeeker map[string]interface{}) (Seeker, error)
	}
)

// TableName table name in db for gorm
func (Seeker) TableName() string {
	return "seekers"
}
