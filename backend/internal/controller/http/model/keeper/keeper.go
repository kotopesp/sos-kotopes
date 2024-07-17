package keeper

type (
	Keeper struct {
		Description string `json:"description"`
		Location    string `json:"location" validate:"required"`
	}
)
