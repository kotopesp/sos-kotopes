package role

type Role struct {
	Name        string `json:"name"`
	Id          int    `gorm:"primary_key;autoIncrement" json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
