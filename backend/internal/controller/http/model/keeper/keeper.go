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
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Location    string  `json:"location"`
}

type KeepersUpdate struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Location    string  `json:"location"`
}

type GetAllKeepersParams struct {
	SortBy    string  `query:"sortby"`
	SortOrder string  `query:"sortorder"`
	Location  string  `query:"location"`
	MinRating float64 `query:"minrating"`
	MaxRating float64 `query:"maxrating"`
	MinPrice  float64 `query:"minprice"`
	MaxPrice  float64 `query:"maxprice"`
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
}
