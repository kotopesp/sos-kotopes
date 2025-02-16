package core

import (
	"context"
)

type (
	Moderator struct {
		UserID int `gorm:"column:user_id"`
	}

	ModeratorStore interface {
		GetModeratorByID(ctx context.Context, id int) (moderator Moderator, err error)
		AddModerator(ctx context.Context, id int) (err error)
	}
	ModeratorService interface {
		GetModerator(ctx context.Context, id int) (moderator Moderator, err error)
	}
)

// TableName table name in db for gorm
func (Moderator) TableName() string {
	return "moderators"
}
