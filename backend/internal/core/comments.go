package core

import (
	"context"
	"time"
)

type Comments struct {
	Id         int       `db:"id"`
	Content    string    `db:"content"`
	Author_id  int       `db:"author_id"`
	Posts_id   int       `db:"posts_id"`
	Is_deleted bool      `db:"is_deleted"`
	Deleted_at time.Time `db:"deleted_at"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Parent_id  int       `db:"parent_id"`
	Reply_id   int       `db:"reply_id"`
}

type CommentsStore interface {
	CreateComment(ctx context.Context, comment Comments, post_id int) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, post_id int) (data []Comments, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, id int) error
}

type CommentsService interface {
	CreateComment(ctx context.Context, comment Comments, post_id int) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, post_id int) (data []Comments, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, id int) error
}

type GetAllParamsComments struct {
	Limit  *int
	Offset *int
}
