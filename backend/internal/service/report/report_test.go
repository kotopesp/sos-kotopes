package report_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/kotopesp/sos-kotopes/internal/core"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/kotopesp/sos-kotopes/internal/service/report"
	"github.com/stretchr/testify/assert"
)

func TestCreateReport_PostAlreadyOnModeration(t *testing.T) {
	// Проверяет, что жалоба не создается, если пост уже на модерации.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	post := core.Post{ID: 1, Status: core.OnModeration}
	mockPosts.On("GetPostByID", ctx, 1).Return(post, nil)

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 1})
	assert.NoError(t, err)
	mockPosts.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "CreateReport", ctx, mock.Anything)
}

func TestCreateReport_Success(t *testing.T) {
	// Проверяет полный путь успешного создания жалобы и отправки на модерацию.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	post := core.Post{ID: 2, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 2).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 2).Return(core.ReportAmountThreshold, nil)
	mockPosts.On("SendToModeration", ctx, 2).Return(nil)

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 2})
	assert.NoError(t, err)
	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestCreateReport_GetPostFails(t *testing.T) {
	// Проверяет, что ошибка получения поста приводит к возврату ошибки.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	mockPosts.On("GetPostByID", ctx, 3).Return(core.Post{}, errors.New("not found"))

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 3})
	assert.Error(t, err)
	mockPosts.AssertExpectations(t)
}

func TestCreateReport_CreateReportFails(t *testing.T) {
	// Проверяет, что ошибка при создании жалобы возвращается корректно.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	post := core.Post{ID: 4, Status: core.Deleted}
	mockPosts.On("GetPostByID", ctx, 4).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(errors.New("create error"))

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 4})
	assert.Error(t, err)
	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestCreateReport_GetReportsCountFails(t *testing.T) {
	// Проверяет, что ошибка при подсчете жалоб обрабатывается корректно.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	post := core.Post{ID: 5, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 5).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 5).Return(0, errors.New("count db error"))

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 5})
	assert.Error(t, err)
	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestCreateReport_SendToModerationFails(t *testing.T) {
	// Проверяет, что ошибка при отправке поста на модерацию корректно возвращается.
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	post := core.Post{ID: 6, Status: core.Published}
	mockPosts.On("GetPostByID", ctx, 6).Return(post, nil)
	mockReports.On("CreateReport", ctx, mock.Anything).Return(nil)
	mockReports.On("GetReportsCount", ctx, 6).Return(core.ReportAmountThreshold, nil)
	mockPosts.On("SendToModeration", ctx, 6).Return(errors.New("send fail"))

	svc := report.NewReportService(mockReports, mockPosts)

	err := svc.CreateReport(ctx, core.Report{PostID: 6})
	assert.Error(t, err)
	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}
