package post

import (
	"time"
)

type (
	Post struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		UserID    int       `json:"user_id"`
		// Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		AnimalID  int       `json:"animal_id"`
		Photo     []byte    `json:"photo"`
	}

	GetAllPostsParams struct {
		Limit      int    `query:"limit"`
		Offset     int    `query:"offset"`
		SortBy     string `query:"sortBy"`
		SortOrder  string `query:"sortOrder"`
		SearchTerm string `query:"q"`
	}
)
