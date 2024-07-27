package comments

import "time"

type Comments struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	PostsID   int       `json:"posts_id" db:"posts_id"`
	IsDeleted bool      `json:"is_deleted" db:"is_deleted"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ParentID  int       `json:"parent_id " db:"parent_id"`
	ReplyID   int       `json:"reply_id" db:"reply_id"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
