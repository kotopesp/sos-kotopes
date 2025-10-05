package comment

import (
	"time"

	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
)

type Comment struct {
	ID        int       `json:"id" example:"3" validate:"required"`
	ParentID  *int      `json:"parent_id" form:"parent_id" example:"1"`
	ReplyID   *int      `json:"reply_id" form:"reply_id" example:"3"`
	User      User      `json:"user" validate:"required"`
	Content   string    `json:"content" form:"content" validate:"required" example:"Hello, world!"`
	Status    string    `json:"status" example:"deleted" oneof:"deleted,published,on_moderation" validate:"required"`
	CreatedAt time.Time `json:"created_at" example:"2021-09-01T12:00:00Z" validate:"required"`
}

type Create struct {
	Content  string `json:"content" form:"content" validate:"required" example:"Hello, world!" minlength:"1"`
	ParentID *int   `json:"parent_id" form:"parent_id" example:"1" min:"1"`
	ReplyID  *int   `json:"reply_id" form:"reply_id" example:"3" min:"1"`
}

type User struct {
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"Jack_Vorobey123"`
}

type GetAllCommentsResponse struct {
	Data []Comment             `json:"comments"`
	Meta pagination.Pagination `json:"meta"`
}

type Update struct {
	Content string `json:"content" form:"content" validate:"required" minlength:"1" example:"Hello, world!"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit" validate:"gt=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type PathParams struct {
	PostID    int `params:"post_id" validate:"gt=0"`
	CommentID int `params:"comment_id" validate:"gt=0"`
}

type PostIDPathParams struct {
	PostID int `params:"post_id" validate:"gt=0"`
}
