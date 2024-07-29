package keeper

import "time"

type Keepers struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Location    string    `gorm:"type:varchar(100)" json:"location"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type GetAllKeepersParams struct {
	SortBy    *string
	SortOrder *string
	Location  *string
	MinRating *float64
	MaxRating *float64
	MinPrice  *float64
	MaxPrice  *float64
	Limit     *int
	Offset    *int
}
