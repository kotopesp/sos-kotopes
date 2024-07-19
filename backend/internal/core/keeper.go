package core

import (
	"context"
	"time"
)

type Keepers struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	Location    string    `gorm:"type:varchar(100)" json:"location"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
}

type GetAllKeepersParams struct {
	SortBy    *string
	SortOrder *string
	Limit     *int
	Offset    *int
}

type KeeperStore interface {
	GetAll(ctx *context.Context, params GetAllKeepersParams) ([]Keepers, error)
	GetByID(ctx *context.Context, id int) (Keepers, error)
	Create(ctx *context.Context, keeper Keepers) error
	DeleteById(ctx *context.Context, id int) error
	UpdateById(ctx *context.Context, id int) error
}

type KeeperService interface {
	GetAll(ctx *context.Context, params GetAllKeepersParams) ([]Keepers, error)
	GetByID(ctx *context.Context, id int) (Keepers, error)
	Create(ctx *context.Context, keeper Keepers) error
	DeleteById(ctx *context.Context, id int) error
	UpdateById(ctx *context.Context, id int) error
}
