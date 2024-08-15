package comment

import (
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/pagination"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	ParentID  *int      `json:"parent_id" form:"parent_id"`
	ReplyID   *int      `json:"reply_id" form:"reply_id"`
	User      User      `json:"user"`
	Content   string    `json:"content" form:"content" validate:"required"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetAllCommentsResponse struct {
	model.Response
	Meta pagination.Pagination `json:"meta"`
}

type Update struct {
	Content string `json:"content" form:"content" validate:"required"`
}

type GetAllCommentsParams struct {
	Limit  int `query:"limit" validate:"gt=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type PathParams struct {
	PostID    int `params:"post_id" validate:"omitempty,gt=0"`
	CommentID int `params:"comment_id" validate:"omitempty,gt=0"`
}
