package core

import (
	"context"
	"time"
)

type (
	Seekers struct {
		ID          int       `gorm:"primaryKey;autoIncrement;column:id"`
		UserID      int       `gorm:"column:user_id"`
		Description string    `gorm:"column:description"`
		Location    string    `gorm:"column:location"`
		EquipmentID int       `gorm:"column:id_equipment"`
		Price       int       `gorm:"column:price"`
		HaveCar     bool      `gorm:"column:have_car"`
		CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP();column:created_at"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
		IsDeleted   bool      `gorm:"column:is_deleted"`
		DeletedAt   time.Time `gorm:"column:deleted_at"`
	}

	UpdateSeekers struct {
		ID          int       `gorm:"primaryKey;autoIncrement;column:id"`
		UserID      int       `gorm:"column:user_id"`
		Description string    `gorm:"column:description"`
		Location    string    `gorm:"column:location"`
		EquipmentID int       `gorm:"column:id_equipment"`
		Price       int       `gorm:"column:price"`
		HaveCar     bool      `gorm:"column:have_car"`
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
		CreateSeeker(ctx context.Context, seeker Seekers) (Seekers, error)
		GetSeeker(ctx context.Context, id int) (Seekers, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeekers) (Seekers, error)
		CreateEquipment(ctx context.Context, equipment Equipment) (int, error)
	}

	SeekersStore interface {
		CreateSeeker(ctx context.Context, seeker Seekers) (Seekers, error)
		GetSeeker(ctx context.Context, id int) (Seekers, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeekers) (Seekers, error)
		CreateEquipment(ctx context.Context, equipment Equipment) (int, error)
	}
)

// TableName table name in db for gorm
func (Seekers) TableName() string {
	return "seekers"
}
