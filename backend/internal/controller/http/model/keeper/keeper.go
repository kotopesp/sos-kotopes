package keeper

import "time"

type Keepers struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	Location    string    `gorm:"type:varchar(100)" json:"location"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
}

type GetAllKeepersParams struct {
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}
