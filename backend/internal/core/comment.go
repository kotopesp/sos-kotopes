package core

import (
	"context"
	"errors"
	"time"
)

type Comment struct {
	ID        int        `gorm:"column:id"`
	ParentID  *int       `gorm:"column:parent_id"`
	ReplyID   *int       `gorm:"column:reply_id"`
	AuthorID  int        `gorm:"column:author_id"`
	PostID    int        `gorm:"column:posts_id"`
	Content   string     `gorm:"column:content"`
	IsDeleted bool       `gorm:"column:is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

type CommentStore interface {
	CreateComment(ctx context.Context, comment Comment) (data Comment, err error)
	GetAllComments(ctx context.Context, params GetAllCommentsParams) (data []Comment, total int, err error)
	UpdateComment(ctx context.Context, comments Comment) (data Comment, err error)
	DeleteComment(ctx context.Context, comments Comment) error
	GetCommentByID(ctx context.Context, commentID int) (data Comment, err error)
}

type CommentService interface {
	CreateComment(ctx context.Context, comment Comment) (data Comment, err error)
	GetAllComments(ctx context.Context, params GetAllCommentsParams) (data []Comment, total int, err error)
	UpdateComment(ctx context.Context, comments Comment) (data Comment, err error)
	DeleteComment(ctx context.Context, comments Comment) error
	GetCommentByID(ctx context.Context, commentID int) (data Comment, err error)
}

type GetAllCommentsParams struct {
	PostID int
	Limit  *int
	Offset *int
}

// errors
var (
	ErrCommentAuthorIDMismatch = errors.New("your user_id and db author_id mismatch")
	ErrCommentPostIDMismatch   = errors.New("your posts_id and db posts_id mismatch")
	ErrNoSuchComment           = errors.New("no such comment")
	ErrCommentIsDeleted        = errors.New("comment is deleted")
)

// TableName table name in db for gorm
func (Comment) TableName() string {
	return "comments"
}
