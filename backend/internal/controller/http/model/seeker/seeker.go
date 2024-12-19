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

type GetAllSeekerParams struct {
	Sort       *string `query:"sort" validate:"omitempty"`
	AnimalType *string `query:"animal_type" validate:"omitempty,oneof=dog cat both"`
	Location   *string `query:"location" validate:"omitempty"`
	Price      *int    `query:"min_price" validate:"omitempty,min=-1"`
	HaveCar    *bool   `query:"have_car" validate:"omitempty,boolean"`
	Limit      *int    `query:"limit" validate:"omitempty,min=1"`
	Offset     *int    `query:"offset" validate:"omitempty,min=0"`
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

type ResponseSeekers struct {
	ResponseSeekers []ResponseSeeker `json:"payload"`
}
