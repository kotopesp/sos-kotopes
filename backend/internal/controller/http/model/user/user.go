package user

type (
	User struct {
		Username    string  `form:"username" validate:"required,max=25,no_specials"`
		Password    string  `form:"password" validate:"required,min=8,max=72,contains_digit,contains_uppercase"`
		Description *string `form:"description"`
		Firstname   *string `form:"firstname"`
		Lastname    *string `form:"lastname"`
		Photo       *[]byte
	}
)
