package vet

type (
	Vet struct {
		Description string `json:"description"`
		Location    string `json:"location" validate:"required"`
	}
)
