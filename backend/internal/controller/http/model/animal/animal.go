package animal

type (
	// UpdateAnimal is the structure used for updating an animal
	UpdateAnimal struct {
		AnimalType  *string `form:"animal_type" json:"animal_type"`
		Age         *int    `form:"age" json:"age"`
		Color       *string `form:"color" json:"color"`
		Gender      *string `form:"gender" json:"gender"`
		Description *string `form:"description" json:"description"`
		Status      *string `form:"status" json:"status"`
	}
)
