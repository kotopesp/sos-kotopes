package keeper

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/user"
)

// KeepersCreate represents the data required to create a new keeper.
type KeepersCreate struct {
	UserID      int     `form:"user_id" validate:"required,min=0"`
	Description string  `form:"description" validate:"required,notblank,max=600"`
	Price       float64 `form:"price" validate:"required,min=0"`
	Location    string  `form:"location" validate:"required"`
}

// KeepersUpdate represents the data to update an existing keeper.
type KeepersUpdate struct {
	ID          int     `form:"id"`
	Description string  `form:"description" validate:"notblank,max=600"`
	Price       float64 `form:"price" validate:"min=0"`
	Location    string  `form:"location"`
}

// KeepersResponse represents the response keeper entity.
type KeepersResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description" validate:"required,notblank,max=600"`
	Price       float64   `json:"price" validate:"min=0"`
	Location    string    `json:"location" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type KeepersResponseWithUser struct {
	Keeper KeepersResponse
	User   user.ResponseUser
}

// KeepersResponseWithMeta represents the respose keeper entity with meta information.
type KeepersResponseWithMeta struct {
	Meta pagination.Pagination     `json:"meta"`
	Data []KeepersResponseWithUser `json:"data"`
}

// GetAllKeepersParams represents the query parameters for filtering and sorting keepers.
type GetAllKeepersParams struct {
	Sort      string   `query:"sort" validate:"omitempty,sort_keeper"`
	Location  *string  `query:"location"`
	MinRating *float64 `query:"min_rating" validate:"omitempty,gte=1,lte=5"`
	MaxRating *float64 `query:"max_rating" validate:"omitempty,gte=1,lte=5"`
	MinPrice  *float64 `query:"min_price" validate:"omitempty,gte=0"`
	MaxPrice  *float64 `query:"max_price" validate:"omitempty,gte=0"`
	Limit     int      `query:"limit" validate:"omitempty,gt=0"`
	Offset    int      `query:"offset" validate:"omitempty,gte=0"`
}
