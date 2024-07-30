package core

import (
	"context"
	"errors"
	"time"
)

type Comments struct {
	ID        int        `gorm:"column:id"`
	Content   string     `gorm:"column:content"`
	AuthorID  int        `gorm:"column:author_id"`
	PostsID   int        `gorm:"column:posts_id"`
	IsDeleted bool       `gorm:"column:is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	ParentID  *int       `gorm:"column:parent_id"`
	ReplyID   *int       `gorm:"column:reply_id"`
}

type CommentsStore interface {
	CreateComment(ctx context.Context, comment Comments) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, postID int) (data []Comments, total int, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, comments Comments) error
	GetCommentByID(ctx context.Context, commentID int) (Comments, error)
}

type CommentsService interface {
	CreateComment(ctx context.Context, comment Comments) (data Comments, err error)
	GetCommentsByPostID(ctx context.Context, params GetAllParamsComments, postID int) (data []Comments, total int, err error)
	UpdateComments(ctx context.Context, comments Comments) (Comments, error)
	DeleteComments(ctx context.Context, comments Comments) error
	GetCommentByID(ctx context.Context, commentID int) (Comments, error)
}

type GetAllParamsComments struct {
	Limit  *int
	Offset *int
}

// errors
var (
	ErrCommentAuthorIDMismatch = errors.New("user_id and comment author_id mismatch")
	ErrNoSuchComment           = errors.New("no such comment")
	ErrNoSuchPost              = errors.New("no such post")
	ErrCommentIsDeleted        = errors.New("comment is deleted")
)
