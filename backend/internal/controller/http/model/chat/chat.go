package conversation

import "time"

type (
	Chat struct {
		ID        int       `json:"id"`
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
	}
)
