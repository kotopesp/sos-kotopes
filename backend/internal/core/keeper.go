package core

import "context"

type Keeper struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int     `json:"user_id"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	Location    string  `gorm:"type:varchar(100)" json:"location"`
}

type GetAllKeepersParams struct {
	SortBy    *string
	SortOrder *string
	Limit     *int
	Offset    *int
}

type KeeperStore interface {
	GetAll(ctx *context.Context, params GetAllKeepersParams) ([]Keeper, error)
	GetByID(ctx *context.Context, id int) (Keeper, error)
	Create(ctx *context.Context, keeper Keeper) error
	DeleteById(ctx *context.Context, id int) error
	UpdateById(ctx *context.Context, id int) error
}

type KeeperService interface {
	GetAll(ctx *context.Context, params GetAllKeepersParams) ([]Keeper, error)
	GetByID(ctx *context.Context, id int) (Keeper, error)
	Create(ctx *context.Context, keeper Keeper) error
	DeleteById(ctx *context.Context, id int) error
	UpdateById(ctx *context.Context, id int) error
}
