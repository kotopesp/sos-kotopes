package post

import "time"

type (
	Post struct {
		ID        int       `json:"id"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		AnimalID  int       `json:"animal_id"`
	}

	GetAllPostsParams struct {
		Limit      int    `query:"limit"`
		Offset     int    `query:"offset"`
		SortBy     string `query:"sort_by"`
		SortOrder  string `query:"sort_order"`
		SearchTerm string `query:"q"`
	}
)
