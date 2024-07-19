package keeper

type Keeper struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int     `json:"user_id"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	Location    string  `gorm:"type:varchar(100)" json:"location"`
}

type GetAllKeepersParams struct {
	SortBy    string
	SortOrder string
	Limit     int
	Offset    int
}
