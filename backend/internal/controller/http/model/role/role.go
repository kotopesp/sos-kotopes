package role

type Role struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GiveRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRole struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
