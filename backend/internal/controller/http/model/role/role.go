package role

import "time"

type Role struct {
	Name        string    `json:"name"`
	ID          int       `json:"id"`
	Username    string    `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GivenRole struct {
	Name        string `json:"name" validate:"required,oneof=keeper seeker vet"`
	Description string `json:"description" validate:"max=512"`
}

type UpdateRole struct {
	Name        string  `json:"name" validate:"required,oneof=keeper seeker vet"`
	Description *string `json:"description" validate:"max=512"`
}

type DeleteRole struct {
	Name string `json:"name" validate:"required,oneof=keeper seeker vet"`
}
