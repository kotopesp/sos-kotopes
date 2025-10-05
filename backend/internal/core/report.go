package core

import (
	"context"
	"time"
)

// Report consists of counter for report amount and reasons why post was reported.
type (
	Report struct {
		ID             int          `gorm:"column:id;primaryKey"`
		UserID         int          `gorm:"column:user_id"`
		Reason         ReportReason `gorm:"column:reason"`
		CreatedAt      time.Time    `gorm:"column:created_at"`
		ReportableID   int          `gorm:"column:reportable_id"`
		ReportableType string       `gorm:"column:reportable_type"`
	}

	ReportStore interface {
		CreateReport(ctx context.Context, report Report) (err error)
		GetReportsCount(ctx context.Context, reportableID int, reportableType string) (int, error)
		GetReportReasons(ctx context.Context, reportableID int, reportableType string) (reasons []string, err error)
		DeleteAllReports(ctx context.Context, reportableID int, reportableType string) (err error)
	}

	ReportService interface {
		CreateReport(ctx context.Context, report Report) (err error)
	}

	// ReportReason is custom type that represents values that can be used for report reasons.
	ReportReason string
)

const (
	Spam           ReportReason = "spam"
	ViolentContent ReportReason = "violent_content"
	ViolentSpeech  ReportReason = "violent_speech"
)

// Status is a custom type that represents the current state of a post.
type ContentStatus string

const (
	Published    ContentStatus = "published"
	Deleted      ContentStatus = "deleted"
	OnModeration ContentStatus = "on_moderation"
)

const (
	ReportableTypePost    = "post"
	ReportableTypeComment = "comment"
)

// ReportAmountThreshold defines the maximum number of reports a post can receive before it is moved to moderation.
const ReportAmountThreshold = 15

func (Report) TableName() string { return "reports" }
