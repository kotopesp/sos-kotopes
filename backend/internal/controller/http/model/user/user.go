package user

type User struct {
	Id        int    `gorm:"primary key;autoIncrement" json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type UpdateUser struct {
	Id       int     `gorm:"primary key;autoIncrement" json:"id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	//CreatedAt    *string `json:"created_at"`
}
