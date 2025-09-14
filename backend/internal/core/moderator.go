package core

import (
	"context"
	"time"
)

type (
	Moderator struct {
		UserID    int       `gorm:"column:user_id"`    // ID of user who is moderator
		CreatedAt time.Time `gorm:"column:created_at"` // TimeStamp shows when this moderator was added
	}

	ModeratorStore interface {
		GetModeratorByID(ctx context.Context, id int) (moderator Moderator, err error)
		CreateModerator(ctx context.Context, moderator Moderator) (err error)
	}
	CommentForModeration struct {
		Comment Comment
		Reasons []string
	}

	ModeratorService interface {
		GetModerator(ctx context.Context, id int) (moderator Moderator, err error)
		GetPostsForModeration(ctx context.Context, filter Filter) (posts []PostForModeration, err error)
		DeletePost(ctx context.Context, ID int) (err error)
		ApprovePost(ctx context.Context, postID int) (err error)
		DeleteComment(ctx context.Context, commentID int) error
		ApproveComment(ctx context.Context, commentID int) error
		GetCommentsForModeration(ctx context.Context, filter Filter) ([]CommentForModeration, error)
		BanUser(ctx context.Context, banRecord BannedUserRecord) error
	}
)

// TableName table name in db for gorm
func (Moderator) TableName() string {
	return "moderators"
}
