package seeker

type CreateSeeker struct {
	AnimalType       string `json:"animal_type" validate:"required,oneof=dog cat both"`
	Description      string `json:"description" validate:"omitempty,max=4000"`
	Location         string `json:"location" validate:"required"`
	EquipmentRental  int    `json:"equipment_rental" validate:"required,min=-1"`
	HaveMetalCage    bool   `json:"have_metal_cage"`
	HavePlasticCage  bool   `json:"have_plastic_cage"`
	HaveNet          bool   `json:"have_net"`
	HaveLadder       bool   `json:"have_ladder"`
	HaveOther        string `json:"have_other"`
	HaveCar          bool   `json:"have_car"`
	Price            int    `json:"price" validate:"min=0"`
	WillingnessCarry string `json:"willingness_carry" validate:"required,oneof=yes no situational"`
}

type UpdateSeeker struct {
	AnimalType       *string `json:"animal_type" validate:"omitempty,oneof=dog cat both"`
	Description      *string `json:"description" validate:"omitempty,max=4000"`
	Location         *string `json:"location" validate:"omitempty"`
	EquipmentRental  *int    `json:"equipment_rental" validate:"omitempty,min=-1"`
	HaveMetalCage    *bool   `json:"have_metal_cage"`
	HavePlasticCage  *bool   `json:"have_plastic_cage"`
	HaveNet          *bool   `json:"have_net"`
	HaveLadder       *bool   `json:"have_ladder"`
	HaveOther        *string `json:"have_other"`
	HaveCar          *bool   `json:"have_car" validate:"omitempty"`
	Price            *int    `json:"price" validate:"omitempty"`
	WillingnessCarry *string `json:"willingness_carry" validate:"omitempty,oneof=yes no situational"`
}

type GetAllSeekerParams struct {
	SortBy             *string `query:"sort_by" validate:"omitempty,oneof=animal_type location price have_car"`
	SortOrder          *string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	AnimalType         *string `query:"animal_type" validate:"omitempty,oneof=dog cat both"`
	Location           *string `query:"location" validate:"omitempty"`
	MinEquipmentRental *int    `query:"min_equipment_rental" validate:"omitempty,min=-1"`
	MaxEquipmentRental *int    `query:"max_equipment_rental" validate:"omitempty,min=-1"`
	HaveMetalCage      *bool   `query:"have_metal_cage"`
	HavePlasticCage    *bool   `query:"have_plastic_cage"`
	HaveNet            *bool   `query:"have_net"`
	HaveLadder         *bool   `query:"have_ladder"`
	HaveOther          *string `query:"have_other"`
	MinPrice           *int    `query:"min_price" validate:"omitempty,min=0"`
	MaxPrice           *int    `query:"max_price" validate:"omitempty,min=0"`
	HaveCar            *bool   `query:"have_car" validate:"omitempty,boolean"`
	Limit              *int    `query:"limit" validate:"omitempty,min=1"`
	Offset             *int    `query:"offset" validate:"omitempty,min=0"`
}

type ResponseSeeker struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	AnimalType       string `json:"animal_type"`
	Description      string `json:"description"`
	Location         string `json:"location"`
	EquipmentRental  int    `json:"equipment_rental"`
	HaveMetalCage    bool   `json:"have_metal_cage"`
	HavePlasticCage  bool   `json:"have_plastic_cage"`
	HaveNet          bool   `json:"have_net"`
	HaveLadder       bool   `json:"have_ladder"`
	HaveOther        string `json:"have_other"`
	HaveCar          bool   `json:"have_car"`
	Price            int    `json:"price"`
	WillingnessCarry string `json:"willingness_carry" validate:"required,oneof=yes no situational"`
}

type ResponseSeekers struct {
	ResponseSeekers []ResponseSeeker `json:"payload"`
}
