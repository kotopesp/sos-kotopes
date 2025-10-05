package report_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kotopesp/sos-kotopes/internal/core"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/kotopesp/sos-kotopes/internal/service/report"
)

func TestCreateReport_PostAlreadyOnModeration(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	post := core.Post{ID: 1, Status: core.OnModeration}
	mockPosts.On("GetPostByID", ctx, 1).Return(post, nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   1,
		ReportableType: core.ReportableTypePost,
	})
	assert.NoError(t, err)

	mockPosts.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "GetReportsCount", ctx, mock.Anything, mock.Anything)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
	mockComments.AssertNotCalled(t, "GetCommentByID", ctx, mock.Anything)
}

func TestCreateReport_CommentAlreadyOnModeration(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	comment := core.Comment{ID: 1, Status: core.OnModeration}
	mockComments.On("GetCommentByID", ctx, 1).Return(comment, nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   1,
		ReportableType: core.ReportableTypeComment,
	})
	assert.NoError(t, err)

	mockComments.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
	mockPosts.AssertNotCalled(t, "GetPostByID", ctx, mock.Anything)
}

func TestCreateReport_PostSuccess(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	post := core.Post{ID: 2, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 2).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 2, core.ReportableTypePost).Return(core.ReportAmountThreshold, nil)
	mockPosts.On("SendToModeration", ctx, 2).Return(nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   2,
		ReportableType: core.ReportableTypePost,
	})
	assert.NoError(t, err)

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
	mockComments.AssertNotCalled(t, "GetCommentByID", ctx, mock.Anything)
}

func TestCreateReport_CommentSuccess(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	comment := core.Comment{ID: 3, Status: core.Published}
	mockComments.On("GetCommentByID", ctx, 3).Return(comment, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 3, core.ReportableTypeComment).Return(core.ReportAmountThreshold, nil)
	mockComments.On("SendToModeration", ctx, 3).Return(nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   3,
		ReportableType: core.ReportableTypeComment,
	})
	assert.NoError(t, err)

	mockComments.AssertExpectations(t)
	mockReports.AssertExpectations(t)
	mockPosts.AssertNotCalled(t, "GetPostByID", ctx, mock.Anything)
}

func TestCreateReport_InvalidReportableType(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	mockComments.On("GetReportsCount")

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   1,
		ReportableType: "invalid_type",
	})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core.ErrInvalidReportableType))

	mockPosts.AssertNotCalled(t, "GetPostByID", ctx, mock.Anything)
	mockComments.AssertNotCalled(t, "GetCommentByID", ctx, mock.Anything)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
}

func TestCreateReport_PostNotFound(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	mockPosts.On("GetPostByID", ctx, 4).Return(core.Post{}, errors.New("not found"))

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   4,
		ReportableType: core.ReportableTypePost,
	})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core.ErrTargetNotFound))

	mockPosts.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
}

func TestCreateReport_CommentNotFound(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	mockComments.On("GetCommentByID", ctx, 5).Return(core.Comment{}, errors.New("not found"))

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   5,
		ReportableType: core.ReportableTypeComment,
	})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, core.ErrTargetNotFound))

	mockComments.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
}

func TestCreatePostReport_DuplicateReport(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	post := core.Post{ID: 6, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 6).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(core.ErrDuplicateReport)
	mockReports.On("GetReportsCount", ctx, 6, core.ReportableTypePost).Return(1, nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   6,
		ReportableType: core.ReportableTypePost,
	})
	assert.NoError(t, err)

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestCreateCommentReport_DuplicateReport(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	comment := core.Comment{ID: 6, Status: core.Published}
	mockComments.On("GetCommentByID", ctx, 6).Return(comment, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(core.ErrDuplicateReport)
	mockReports.On("GetReportsCount", ctx, 6, core.ReportableTypeComment).Return(1, nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   6,
		ReportableType: core.ReportableTypeComment,
	})
	assert.NoError(t, err)

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestCreateReport_BelowThreshold(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)
	mockComments := new(mocks.MockCommentStore)

	post := core.Post{ID: 7, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 7).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 7, core.ReportableTypePost).Return(5, nil)

	svc := report.NewReportService(mockReports, mockPosts, mockComments)

	err := svc.CreateReport(ctx, core.Report{
		ReportableID:   7,
		ReportableType: core.ReportableTypePost,
	})
	assert.NoError(t, err)

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
	mockPosts.AssertNotCalled(t, "SendToModeration", ctx, mock.Anything)
}
