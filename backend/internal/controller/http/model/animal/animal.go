package animal

type (
	Animal struct {
		AnimalType  string    `form:"animal_type" json:"animal_type" validate:"required"` 
		Age         int       `form:"age" json:"age" validate:"gte=0"`
		Color       string    `form:"color" json:"color" validate:"required"`
		Gender      string    `form:"gender" json:"gender" validate:"required"`
		Description string    `form:"description" json:"description"`
		Status      string    `form:"status" json:"status" validate:"required"`
	}

	AnimalResponse struct {
		AnimalType  string    `form:"animal_type" json:"animal_type"` 
		Age         int       `form:"age" json:"age"`
		Color       string    `form:"color" json:"color"`
		Gender      string    `form:"gender" json:"gender"`
		Description string    `form:"description" json:"description"`
		Status      string    `form:"status" json:"status"`
	}

	UpdateRequestBodyAnimal struct {
		AnimalType  *string `form:"animal_type" json:"animal_type"`
		Age         *int    `form:"age" json:"age"`
		Color       *string `form:"color" json:"color"`
		Gender      *string `form:"gender" json:"gender"`
		Description *string `form:"description" json:"description"`
		Status      *string `form:"status" json:"status"`
	}
)
