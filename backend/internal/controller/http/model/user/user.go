package user

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Description  string `json:"description"`
	Photo        string `json:"photo"`
	PasswordHash string `json:"password"`
	IsDeleted    bool   `json:"is_deleted"`
	DeletedAt    string `json:"deleted_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UpdateUser struct {
	Id           int     `json:"id"`
	Username     *string `json:"username"`
	FirstName    *string `json:"firstname"`
	LastName     *string `json:"lastname"`
	Description  *string `json:"description"`
	Photo        *string `json:"photo"`
	PasswordHash *string `json:"password"`
	IsDeleted    *bool   `json:"is_deleted"`
	DeletedAt    *string `json:"deleted_at"`
	CreatedAt    *string `json:"created_at"`
}

type GetFavourites struct {
	Count  *int    `json:"count"`
	Offset *int    `json:"offset"`
	Sort   *string `json:"sort"`
}
