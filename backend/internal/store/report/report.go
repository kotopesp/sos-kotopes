package report

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
)

type store struct {
	*postgres.Postgres
}

func New(db *postgres.Postgres) core.ReportStore {
	return &store{db}
}

// CreateReport - creates report record for post.
func (s *store) CreateReport(ctx context.Context, report core.Report) (err error) {
	report.CreatedAt = time.Now().UTC()

	if err = s.DB.WithContext(ctx).Create(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Log().Error(ctx, err.Error())
			return core.ErrDuplicateReport
		}

		logger.Log().Error(ctx, err.Error())

		return core.ErrToCreateReport
	}

	return nil
}

// GetReportsCount - returns number of reports for current post by its ID.
func (s *store) GetReportsCount(ctx context.Context, postID int) (int, error) {
	var reportCount int64
	if err := s.DB.WithContext(ctx).Model(&core.Report{}).
		Where("post_id = ?", postID).
		Count(&reportCount).Error; err != nil {
		logger.Log().Error(ctx, err.Error())

		return 0, err
	}

	return int(reportCount), nil
}

// GetReportReasonsForPost returns list of reasons why post was banned.
func (s *store) GetReportReasonsForPost(ctx context.Context, postID int) (reasons []string, err error) {
	err = s.DB.WithContext(ctx).
		Table(core.Report{}.TableName()).
		Where("post_id = ?", postID).
		Pluck("reason", &reasons).Error
	if err != nil {
		logger.Log().Debug(ctx, err.Error())

		return nil, core.ErrGettingReportReasons
	}

	return reasons, nil
}
