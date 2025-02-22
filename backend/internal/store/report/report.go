package report

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"time"
)

type store struct {
	*postgres.Postgres
}

func New(db *postgres.Postgres) core.ReportStore {
	return &store{db}
}

// CreateReport - creates report record for post.
func (s *store) CreateReport(ctx context.Context, report core.Report) (int, error) {
	report.ReportedAt = time.Now().UTC()

	err := s.DB.WithContext(ctx).Create(&report).Error
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return 0, err
	}

	var reportCount int64
	err = s.DB.WithContext(ctx).
		Model(&core.Report{}).
		Where("post_id = ?", report.PostID).
		Count(&reportCount).Error
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return 0, err
	}

	return int(reportCount), nil
}
