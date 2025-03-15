package report

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"github.com/kotopesp/sos-kotopes/pkg/postgres"
	"gorm.io/gorm"
	"time"
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
