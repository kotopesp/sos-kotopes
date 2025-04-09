package moderator_test

import (
	"context"
	"errors"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"testing"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/service/moderator"

	"github.com/stretchr/testify/assert"
)

func TestGetModerator_Success(t *testing.T) {
	ctx := context.TODO()
	mockMod := new(mocks.MockModeratorStore)
	svc := moderator.New(mockMod, nil, nil)

	expected := core.Moderator{UserID: 1}
	mockMod.On("GetModeratorByID", ctx, 1).Return(expected, nil)

	m, err := svc.GetModerator(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, m)
	mockMod.AssertExpectations(t)
}

func TestGetModerator_Failure(t *testing.T) {
	ctx := context.TODO()
	mockMod := new(mocks.MockModeratorStore)
	svc := moderator.New(mockMod, nil, nil)

	mockMod.On("GetModeratorByID", ctx, 2).Return(core.Moderator{}, core.ErrNoSuchModerator)

	_, err := svc.GetModerator(ctx, 2)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchModerator, err)
	mockMod.AssertExpectations(t)
}

// Testing that some
func TestGetPostsForModeration_Success(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	filter := core.FilterDESC
	posts := []core.Post{{ID: 1}, {ID: 2}}
	mockPosts.On("GetPostsForModeration", ctx, filter).Return(posts, nil)
	mockReports.On("GetReportReasonsForPost", ctx, 1).Return([]string{"spam"}, nil)
	mockReports.On("GetReportReasonsForPost", ctx, 2).Return([]string{"offensive"}, nil)

	svc := moderator.New(nil, mockPosts, mockReports)

	result, err := svc.GetPostsForModeration(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "spam", result[0].Reasons[0])
	assert.Equal(t, "offensive", result[1].Reasons[0])

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

// Testing that after some error happened extracting report reasons list, extraction continues.
func TestGetPostsForModeration_ReportFail_Continues(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	filter := core.FilterASC
	posts := []core.Post{{ID: 1}, {ID: 2}}
	mockPosts.On("GetPostsForModeration", ctx, filter).Return(posts, nil)
	mockReports.On("GetReportReasonsForPost", ctx, 1).Return(nil, core.ErrGettingReportReasons)
	mockReports.On("GetReportReasonsForPost", ctx, 2).Return([]string{"spam"}, nil)

	svc := moderator.New(nil, mockPosts, mockReports)

	result, err := svc.GetPostsForModeration(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 2, result[0].Post.ID)
	assert.Equal(t, "spam", result[0].Reasons[0])

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestGetPostsForModeration_NoPostsForModeration(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	filter := core.Filter("asc")
	mockPosts.On("GetPostsForModeration", ctx, filter).Return(nil, core.ErrNoPostsWaitingForModeration)

	svc := moderator.New(nil, mockPosts, mockReports)

	listOfPosts, err := svc.GetPostsForModeration(ctx, filter)
	assert.Error(t, err)
	assert.Equal(t, err, core.ErrNoPostsWaitingForModeration)
	assert.Empty(t, listOfPosts)
	mockPosts.AssertExpectations(t)
}

func TestDeletePost_Success(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockPosts.On("DeletePost", ctx, 10).Return(nil)

	svc := moderator.New(nil, mockPosts, nil)
	err := svc.DeletePost(ctx, 10)

	assert.NoError(t, err)
	mockPosts.AssertExpectations(t)
}

func TestDeletePost_Failure(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockPosts.On("DeletePost", ctx, 99).Return(core.ErrPostNotFound)

	svc := moderator.New(nil, mockPosts, nil)
	err := svc.DeletePost(ctx, 99)

	assert.Error(t, err)
	assert.Equal(t, err, core.ErrPostNotFound)
	mockPosts.AssertExpectations(t)
}

func TestApprovePost_Success(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)
	mockReports.On("DeleteAllReportsForPost", ctx, 5).Return(nil)
	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(nil)

	svc := moderator.New(nil, mockPosts, mockReports)
	err := svc.ApprovePost(ctx, 5)

	assert.NoError(t, err)
	mockReports.AssertExpectations(t)
}

func TestApprovePost_Failure(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)
	mockReports.On("DeleteAllReportsForPost", ctx, 5).Return(errors.New("fail"))
	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(nil)

	svc := moderator.New(nil, mockPosts, mockReports)
	err := svc.ApprovePost(ctx, 5)

	assert.Error(t, err)
	mockReports.AssertExpectations(t)
}

func TestApprovePost_Failure_ApprovePostFromModeration(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)

	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(errors.New("approve failed"))

	svc := moderator.New(nil, mockPosts, mockReports)
	err := svc.ApprovePost(ctx, 5)

	assert.Error(t, err)
	assert.EqualError(t, err, "approve failed")
	mockPosts.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "DeleteAllReportsForPost", ctx, 5)
}
