package role

import "time"

type Role struct {
	Name        string    `json:"name"`
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GiveRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRole struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
