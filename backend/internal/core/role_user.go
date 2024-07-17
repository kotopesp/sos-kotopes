package core

type (
	RoleUser struct {
		ID     int `gorm:"column:id"`
		RoleID int `gorm:"column:role_id"`
		UserID int `gorm:"column:user_id"`
	}
)

// table name in db for gorm
func (RoleUser) TableName() string {
	return "roles_users"
}
