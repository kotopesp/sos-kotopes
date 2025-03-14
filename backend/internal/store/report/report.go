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
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// here i checking if posts exists
	var postExists int64
	if err := tx.Model(&core.Post{}).
		Where("id = ? AND status = ?", report.PostID, string(core.Published)).
		Count(&postExists).Error; err != nil {
		logger.Log().Error(ctx, err.Error())

		return 0, err
	}

	if postExists == 0 {
		logger.Log().Error(ctx, core.ErrPostNotFound.Error())

		return 0, core.ErrPostNotFound
	}

	report.CreatedAt = time.Now().UTC()
	if err := tx.Create(&report).Error; err != nil {
		logger.Log().Error(ctx, err.Error())

		return 0, core.ErrToCreateReport
	}

	var reportCount int64
	if err := tx.Model(&core.Report{}).
		Where("post_id = ?", report.PostID).
		Count(&reportCount).Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Log().Error(ctx, err.Error())
		return 0, err
	}

	return int(reportCount), nil
}
