package chat

import "time"

type (
	Chat struct {
		ID        int       `json:"id"`
		ChatType  string    `json:"chat_type"`
		IsDeleted bool      `json:"is_deleted"`
		DeletedAt time.Time `json:"deleted_at"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
