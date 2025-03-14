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
		GetPostsForModeration(ctx context.Context) (posts []PostForModeration, err error)
	}

	ModeratorService interface {
		GetModerator(ctx context.Context, id int) (moderator Moderator, err error)
		GetPostsForModeration(ctx context.Context) (posts []PostForModeration, err error)
	}
)

// TableName table name in db for gorm
func (Moderator) TableName() string {
	return "moderators"
}
