package seeker

type CreateSeeker struct {
	UserID      int      `json:"user_id" validate:"required,min=1"`
	Description string   `json:"description" validate:"omitempty,max=4000"`
	Location    string   `json:"location" validate:"required"`
	Equipment   []string `json:"equipment" validate:"required"`
	HaveCar     bool     `json:"have_car"`
	Price       int      `json:"price" validate:"min=0"`
}

type UpdateSeeker struct {
	UserID      *int      `json:"user_id" validate:"required,min=1"`
	Description *string   `json:"description" validate:"omitempty,max=4000"`
	Location    *string   `json:"location" validate:"omitempty"`
	Equipment   *[]string `json:"equipment" validate:"omitempty"`
	HaveCar     *bool     `json:"have_car" validate:"omitempty"`
	Price       *int      `json:"price" validate:"omitempty"`
}

type ResponseSeeker struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Equipment   int    `json:"equipment"`
	HaveCar     bool   `json:"have_car"`
	Price       int    `json:"price"`
}
