package response

import "time"

type Response struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	AuthorID  int       `json:"author_id"`
	Content   string    `json:"content"`
	IsDeleted bool      `json:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllResponsesParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
