package user

type (
	User struct {
		Username    string `json:"username" validate:"required"`
		Password    string `json:"password" validate:"required,min=8,max=72,containsDigit,containsUppercase"`
		Description string `json:"description"`
	}
)
