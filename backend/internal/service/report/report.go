package report

import (
	"context"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

type service struct {
	reportStore core.ReportStore
	postStore   core.PostStore
}

func NewReportService(reportStore core.ReportStore, postStore core.PostStore) core.ReportService {
	return &service{reportStore: reportStore, postStore: postStore}
}

// CreateReport - creates Report record in special table.
func (s *service) CreateReport(ctx context.Context, report core.Report) (err error) {
	post, err := s.postStore.GetPostByID(ctx, report.PostID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	// post already on moderation
	if post.Status == core.OnModeration {
		return nil
	}

	if err = s.reportStore.CreateReport(ctx, report); err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	reportCount, err := s.reportStore.GetReportsCount(ctx, report.PostID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	if reportCount >= core.ReportAmountThreshold {
		err = s.postStore.SendToModeration(ctx, report.PostID)
		if err != nil {
			return err
		}
	}

	return nil
}
