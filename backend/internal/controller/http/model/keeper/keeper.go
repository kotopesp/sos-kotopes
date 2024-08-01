package keeper

import "time"

type Keepers struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type KeepersCreate struct {
	UserID      int     `json:"user_id"`
	Description string  `json:"description" validate:"required,notblank,max=600"`
	Price       float64 `json:"price" validate:"required, min=0"`
	Location    string  `json:"location"`
}

type KeepersUpdate struct {
	Description string  `json:"description" validate:"notblank,max=600"`
	Price       float64 `json:"price" validate:"min=0"`
	Location    string  `json:"location"`
}

type GetAllKeepersParams struct {
	Sort      string  `query:"sort"`
	Location  string  `query:"location"`
	MinRating float64 `query:"min_rating"`
	MaxRating float64 `query:"max_rating"`
	MinPrice  float64 `query:"min_price"`
	MaxPrice  float64 `query:"max_price"`
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
}
