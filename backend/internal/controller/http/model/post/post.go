package post

import (
	"time"
)

type (
	Post struct {
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Body      string `json:"body"`
		UserID    int    `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		AnimalID  int    `json:"animal_id"`
		Photo     []byte `json:"photo"`
	}
)