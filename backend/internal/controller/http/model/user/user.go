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
)
