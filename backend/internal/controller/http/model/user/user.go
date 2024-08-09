package user

type (
	User struct {
		Username    string  `form:"username" validate:"required,max=50,no_specials"`
		Password    string  `form:"password" validate:"required,min=8,max=72,contains_digit,contains_uppercase"`
		Description *string `form:"description" validate:"omitempty,max=4000"`
		Firstname   *string `form:"firstname" validate:"omitempty,max=25"`
		Lastname    *string `form:"lastname" validate:"omitempty,max=25"`
		Photo       *[]byte
	}
	UpdateUser struct {
		Username    *string `json:"username" validate:"omitempty,max=50,no_specials"`
		Firstname   *string `json:"firstname" validate:"omitempty,max=25"`
		Lastname    *string `json:"lastname" validate:"omitempty,max=25"`
		Description *string `json:"description" validate:"omitempty,max=512"`
		Photo       *[]byte `json:"photo"`
		Password    *string `json:"password" validate:"omitempty,min=8,max=72,contains_digit,contains_uppercase"`
	}

	ResponseUser struct {
		ID          int     `json:"id"`
		Username    string  `json:"username"`
		Firstname   *string `json:"firstname"`
		Lastname    *string `json:"lastname"`
		Description *string `json:"description"`
		Photo       *[]byte `json:"photo"`
	}

	GetFavourites struct {
		Count  *int    `json:"count"`
		Offset *int    `json:"offset"`
		Sort   *string `json:"sort"`
	}
)
