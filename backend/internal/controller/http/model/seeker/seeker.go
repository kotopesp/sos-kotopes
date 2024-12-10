package seeker

type ResponseSeekers struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Equipment   string `json:"equipment"`
	HaveCar     bool   `json:"have_car"`
	Price       int    `json:"price"`
}
