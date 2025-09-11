package core

import (
	"context"
	"time"
)

type Comment struct {
	ID        int           `gorm:"column:id" fake:"{number:1,100}"`
	ParentID  *int          `gorm:"column:parent_id" fake:"{number:1,100}"`
	ReplyID   *int          `gorm:"column:reply_id" fake:"{number:1,100}"`
	PostID    int           `gorm:"column:posts_id" fake:"{number:1,100}"`
	Status    ContentStatus `gorm:"column:status"`
	AuthorID  int           `gorm:"column:author_id" fake:"{number:1,100}"`
	Author    User          `gorm:"foreignKey:AuthorID;references:ID" fake:"skip"`
	Content   string        `gorm:"column:content" fake:"{sentence:3}"`
	DeletedAt time.Time     `gorm:"column:deleted_at" fake:"skip"`
	CreatedAt time.Time     `gorm:"column:created_at" fake:"skip"`
	UpdatedAt time.Time     `gorm:"column:updated_at" fake:"skip"`
}

type CommentStore interface {
	CreateComment(ctx context.Context, comment Comment) (data Comment, err error)
	GetAllComments(ctx context.Context, params GetAllCommentsParams) (data []Comment, total int, err error)
	UpdateComment(ctx context.Context, comments Comment) (data Comment, err error)
	DeleteComment(ctx context.Context, comments Comment) error
	GetCommentByID(ctx context.Context, commentID int) (data Comment, err error)
	SendToModeration(ctx context.Context, commentID int) error
	GetCommentsForModeration(ctx context.Context, filter Filter) ([]Comment, error)
	ApproveCommentFromModeration(ctx context.Context, commentID int) error
}

type CommentService interface {
	CreateComment(ctx context.Context, comment Comment) (data Comment, err error)
	GetAllComments(ctx context.Context, params GetAllCommentsParams) (data []Comment, total int, err error)
	UpdateComment(ctx context.Context, comments Comment) (data Comment, err error)
	DeleteComment(ctx context.Context, comments Comment) error
}

type GetAllCommentsParams struct {
	PostID int
	Limit  *int
	Offset *int
}

const AmountOfCommentsForModeration = 10

// TableName table name in db for gorm
func (Comment) TableName() string {
	return "comments"
}
