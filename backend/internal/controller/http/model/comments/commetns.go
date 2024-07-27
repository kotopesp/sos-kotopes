package comments

import "time"

type Comments struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"autorID" db:"author_id"`
	PostsID   int       `json:"postsID" db:"posts_id"`
	IsDeleted bool      `json:"isDeleted" db:"is_deleted"`
	DeletedAt time.Time `json:"deletedAt" db:"deleted_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	ParentID  int       `json:"parentID " db:"parent_id"`
	ReplyID   int       `json:"replyID" db:"reply_id"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
