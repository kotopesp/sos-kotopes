package core

type (
	Keeper struct {
		ID          int    `gorm:"column:id"`
		UserID      int    `gorm:"column:user_id"`
		Description string `gorm:"column:description"`
		Location    string `gorm:"column:location"`
	}
)

// table name in db for gorm
func (Keeper) TableName() string {
	return "keepers"
}
