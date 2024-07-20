package user

type (
	User struct {
		Username    string  `json:"username" validate:"required"`
		Password    string  `json:"password" validate:"required,min=8,max=72,contains_digit,contains_uppercase"`
		Description *string `json:"description"`
		Firstname   *string `json:"firstname"`
		Lastname    *string `json:"lastname"`
		Photo       *[]byte `json:"photo"`
	}
)
