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
func (s *store) GetReportsCount(ctx context.Context, reportableID int, reportableType string) (int, error) {
	var reportCount int64
	if err := s.DB.WithContext(ctx).
		Model(&core.Report{}).
		Where("reportable_id = ? AND reportable_type = ?", reportableID, reportableType).
		Count(&reportCount).Error; err != nil {
		logger.Log().Error(ctx, err.Error())

		return 0, err
	}

	return int(reportCount), nil
}

// GetReportReasons returns list of reasons why reportable entity was reported.
func (s *store) GetReportReasons(ctx context.Context, reportableID int, reportableType string) (reasons []string, err error) {
	err = s.DB.WithContext(ctx).
		Model(&core.Report{}).
		Where("reportable_id = ? AND reportable_type = ?", reportableID, reportableType).
		Pluck("DISTINCT reason", &reasons).Error
	if err != nil {
		logger.Log().Debug(ctx, "Failed to get report reasons: "+err.Error())
		return nil, core.ErrGettingReportReasons
	}

	return reasons, nil
}

// DeleteAllReports - delete all report records for specific reportable entity.
func (s *store) DeleteAllReports(ctx context.Context, reportableID int, reportableType string) error {
	if err := s.DB.WithContext(ctx).
		Where("reportable_id = ? AND reportable_type = ?", reportableID, reportableType).
		Delete(&core.Report{}).Error; err != nil {
		logger.Log().Error(ctx, "Failed to delete reports: "+err.Error())
		return core.ErrDeleteReports
	}

	return nil
}
