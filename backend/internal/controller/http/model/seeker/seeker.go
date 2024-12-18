package seeker

type CreateSeeker struct {
	UserID           int      `json:"user_id" validate:"required,min=1"`
	AnimalType       string   `json:"animal_type" validate:"required,oneof=dog cat both"`
	Description      string   `json:"description" validate:"omitempty,max=4000"`
	Location         string   `json:"location" validate:"required"`
	EquipmentRental  int      `json:"equipment_rental" validate:"required,min=-1"`
	Equipment        []string `json:"equipment" validate:"required"`
	HaveCar          bool     `json:"have_car"`
	Price            int      `json:"price" validate:"min=0"`
	WillingnessCarry string   `json:"willingness_carry" validate:"required,oneof=yes no situational"`
}

type UpdateSeeker struct {
	UserID           *int      `json:"user_id" validate:"required,min=1"`
	AnimalType       *string   `json:"animal_type" validate:"omitempty,oneof=dog cat both"`
	Description      *string   `json:"description" validate:"omitempty,max=4000"`
	Location         *string   `json:"location" validate:"omitempty"`
	EquipmentRental  *int      `json:"equipment_rental" validate:"omitempty,min=-1"`
	Equipment        *[]string `json:"equipment" validate:"omitempty"`
	HaveCar          *bool     `json:"have_car" validate:"omitempty"`
	Price            *int      `json:"price" validate:"omitempty"`
	WillingnessCarry *string   `json:"willingness_carry" validate:"omitempty,oneof=yes no situational"`
}

type ResponseSeeker struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	AnimalType       string `json:"animal_type"`
	Description      string `json:"description"`
	Location         string `json:"location"`
	EquipmentRental  int    `json:"equipment_rental"`
	Equipment        int    `json:"equipment"`
	HaveCar          bool   `json:"have_car"`
	Price            int    `json:"price"`
	WillingnessCarry string `json:"willingness_carry" validate:"required,oneof=yes no situational"`
}
