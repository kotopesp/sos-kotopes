package core

import (
	"context"
	"time"
)

type Comments struct {
	ID        int       `db:"id"`
	Content   string    `db:"content"`
	AuthorID  int       `db:"author_id"`
	PostsID   int       `db:"posts_id"`
	IsDeleted bool      `db:"is_deleted"`
	DeletedAt time.Time `db:"deleted_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ParentID  int       `db:"parent_id"`
	ReplyID   int       `db:"reply_id"`
}

type CommentsStore interface {
	CreateComment(ctx context.Context, comment Comments, postID int) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, postID int) (data []Comments, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, id int) error
}

type CommentsService interface {
	CreateComment(ctx context.Context, comment Comments, postID int) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, postID int) (data []Comments, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, id int) error
}

type GetAllParamsComments struct {
	Limit  *int
	Offset *int
}
