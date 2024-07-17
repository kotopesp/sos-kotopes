package seeker

type (
	Seeker struct {
		Description string `json:"description"`
		Location    string `json:"location" validate:"required"`
	}
)
