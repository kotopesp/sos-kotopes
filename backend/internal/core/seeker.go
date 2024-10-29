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
		Equipment   string    `gorm:"column:equipment"`
		Price       int       `gorm:"column:price"`
		Car         bool      `gorm:"column:car"`
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
		Equipment   string    `gorm:"column:equipment"`
		Price       int       `gorm:"column:price"`
		Car         bool      `gorm:"column:car"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime;column:updated_at"`
	}

	SeekersService interface {
		CreateSeeker(ctx context.Context, seeker Seekers) (Seekers, error)
		GetSeeker(ctx context.Context, id int) (Seekers, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeekers) (Seekers, error)
	}

	SeekersStore interface {
		CreateSeeker(ctx context.Context, seeker Seekers) (Seekers, error)
		GetSeeker(ctx context.Context, id int) (Seekers, error)
		UpdateSeeker(ctx context.Context, seeker UpdateSeekers) (Seekers, error)
	}
)

// TableName table name in db for gorm
func (Seekers) TableName() string {
	return "seekers"
}
