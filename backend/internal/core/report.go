package core

import (
	"context"
	"time"
)

// Report consists of counter for report amount and reasons why post was reported.
type (
	Report struct {
		ID         int       `gorm:"column:id;primaryKey"`
		PostID     int       `gorm:"column:post_id"`
		UserID     int       `gorm:"column:user_id"`
		Reason     string    `gorm:"column:reason"`
		ReportedAt time.Time `gorm:"column:reported_at"`
	}

	ReportStore interface {
		CreateReport(ctx context.Context, report Report) (reportCount int, err error)
	}

	ReportService interface {
		CreateReport(ctx context.Context, report Report) (err error)
	}
)

type PostStatus string

const (
	Published    PostStatus = "published"
	Deleted      PostStatus = "deleted"
	OnModeration PostStatus = "on_moderation"
)

const ReportAmountThreshold = 15

func (Report) TableName() string { return "reports" }
