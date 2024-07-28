package user



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

type (
	User struct {
		Username    string  `form:"username" validate:"required,max=50,no_specials"`
		Password    string  `form:"password" validate:"required,min=8,max=72,contains_digit,contains_uppercase"`
		Description *string `form:"description"`
		Firstname   *string `form:"firstname" validate:"omitempty,max=25"`
		Lastname    *string `form:"lastname" validate:"omitempty,max=25"`
		Photo       *[]byte
	}
)

