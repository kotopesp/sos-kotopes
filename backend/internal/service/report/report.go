package report

import (
	"context"
	"errors"

	"github.com/kotopesp/sos-kotopes/internal/core"
)

type service struct {
	reportStore  core.ReportStore
	postStore    core.PostStore
	commentStore core.CommentStore
}

func NewReportService(reportStore core.ReportStore, postStore core.PostStore, commentStore core.CommentStore) core.ReportService {
	return &service{reportStore: reportStore, postStore: postStore, commentStore: commentStore}
}

func (s *service) CreateReport(ctx context.Context, report core.Report) error {
	if err := s.validateTarget(ctx, report); err != nil {
		if errors.Is(err, core.ErrContentAlreadyOnModeration) {
			return nil
		}

		return err
	}

	if err := s.checkModerationThreshold(ctx, report); err != nil {
		return err
	}

	if err := s.createReportRecord(ctx, report); err != nil {
		return err
	}

	return nil
}

func (s *service) validateTarget(ctx context.Context, report core.Report) error {
	switch report.ReportableType {
	case core.ReportableTypePost:
		post, err := s.postStore.GetPostByID(ctx, report.ReportableID)
		if err != nil {
			return core.ErrTargetNotFound
		}
		if post.Status == core.OnModeration {
			return core.ErrContentAlreadyOnModeration
		}

	case core.ReportableTypeComment:
		comment, err := s.commentStore.GetCommentByID(ctx, report.ReportableID)
		if err != nil {
			return core.ErrTargetNotFound
		}
		if comment.Status == core.OnModeration {
			return core.ErrContentAlreadyOnModeration
		}
	default:
		return core.ErrInvalidReportableType
	}

	return nil
}

func (s *service) createReportRecord(ctx context.Context, report core.Report) error {
	err := s.reportStore.CreateReport(ctx, report)
	if errors.Is(err, core.ErrDuplicateReport) {
		return nil
	}

	return err
}

func (s *service) checkModerationThreshold(ctx context.Context, report core.Report) error {
	reportCount, err := s.reportStore.GetReportsCount(ctx, report.ReportableID, report.ReportableType)
	if err != nil {
		return err
	}

	if reportCount >= core.ReportAmountThreshold {
		return s.sendToModeration(ctx, report)
	}

	return nil
}

func (s *service) sendToModeration(ctx context.Context, report core.Report) error {
	switch report.ReportableType {
	case core.ReportableTypePost:
		return s.postStore.SendToModeration(ctx, report.ReportableID)
	case core.ReportableTypeComment:
		return s.commentStore.SendToModeration(ctx, report.ReportableID)
	default:
		return core.ErrInvalidReportableType
	}
}
