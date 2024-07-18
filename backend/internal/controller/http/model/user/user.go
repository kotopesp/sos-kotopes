package user

type User struct {
	Id           int    `gorm:"primary key;autoIncrement" json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
}

type UpdateUser struct {
	Id           int     `gorm:"primary key;autoIncrement" json:"id"`
	Username     *string `json:"username"`
	PasswordHash *string `json:"password_hash"`
	//CreatedAt    *string `json:"created_at"`
}
