package core

import (
	"context"
	"time"
)

// Report consists of counter for report amount and reasons why post was reported.
type (
	Report struct {
		ID        int          `gorm:"column:id;primaryKey"`
		PostID    int          `gorm:"column:post_id"`
		UserID    int          `gorm:"column:user_id"`
		Reason    ReportReason `gorm:"column:reason"`
		CreatedAt time.Time    `gorm:"column:created_at"`
	}

	ReportStore interface {
		CreateReport(ctx context.Context, report Report) (err error)
		GetReportsCount(ctx context.Context, postID int) (int, error)
		GetReportReasonsForPost(ctx context.Context, postID int) ([]string, error)
		DeleteAllReportsForPost(ctx context.Context, postID int) (err error)
	}

	ReportService interface {
		CreateReport(ctx context.Context, report Report) (err error)
	}
)

// ReportReason is custom type that represents values that can be used for report reasons.
type ReportReason string

const (
	Spam           ReportReason = "spam"
	ViolentContent ReportReason = "violent_content"
	ViolentSpeech  ReportReason = "violent_speech"
)

// ReportAmountThreshold defines the maximum number of reports a post can receive before it is moved to moderation.
const ReportAmountThreshold = 15

func (Report) TableName() string { return "reports" }
