package comments

import "time"

type Comments struct {
	ID        int       `json:"id" form:"id"`
	Content   string    `json:"content" form:"content" validate:"required"`
	ParentID  *int      `json:"parent_id" form:"parent_id"`
	ReplyID   *int      `json:"reply_id" form:"reply_id"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted"`
}

type CommentUpdate struct {
	Content string `json:"content" form:"content" validate:"required"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit" validate:"gt=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type CommentURLParams struct {
	PostID    int `params:"post_id" validate:"omitempty,gt=0"`
	CommentID int `params:"comment_id" validate:"omitempty,gt=0"`
}
