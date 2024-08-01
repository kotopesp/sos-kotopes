package keeper

import "time"

// Keepers represents the keeper entity.
type Keepers struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description" validate:"required,notblank,max=600"`
	Price       float64   `json:"price" validate:"min=0"`
	Location    string    `json:"location" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// KeepersCreate represents the data required to create a new keeper.
type KeepersCreate struct {
	UserID      int     `json:"user_id" validate:"required,min=0"`
	Description string  `json:"description" validate:"required,notblank,max=600"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Location    string  `json:"location" validate:"required"`
}

// KeepersUpdate represents the data to update an existing keeper.
type KeepersUpdate struct {
	Description string  `json:"description" validate:"notblank,max=600"`
	Price       float64 `json:"price" validate:"min=0"`
	Location    string  `json:"location"`
}

// GetAllKeepersParams represents the query parameters for filtering and sorting keepers.
type GetAllKeepersParams struct {
	Sort      string  `query:"sort" validate:"omitempty"`
	Location  string  `query:"location"`
	MinRating float64 `query:"min_rating" validate:"omitempty,gte=1,lte=5"`
	MaxRating float64 `query:"max_rating" validate:"omitempty,gte=1,lte=5"`
	MinPrice  float64 `query:"min_price" validate:"omitempty,gte=0"`
	MaxPrice  float64 `query:"max_price" validate:"omitempty,gte=0"`
	Limit     int     `query:"limit" validate:"omitempty,gt=0"`
	Offset    int     `query:"offset" validate:"omitempty,gte=0"`
}
