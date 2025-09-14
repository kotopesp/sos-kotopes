package moderator_test

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"

	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/internal/service/moderator"

	"github.com/stretchr/testify/assert"
)

func TestGetModerator_Success(t *testing.T) {
	ctx := context.TODO()
	mockMod := new(mocks.MockModeratorStore)
	svc := moderator.New(mockMod, nil, nil, nil)

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
	svc := moderator.New(mockMod, nil, nil, nil)

	mockMod.On("GetModeratorByID", ctx, 2).Return(core.Moderator{}, core.ErrNoSuchModerator)

	_, err := svc.GetModerator(ctx, 2)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchModerator, err)
	mockMod.AssertExpectations(t)
}

func TestGetPostsForModeration_Success(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	filter := core.FilterDESC
	posts := []core.Post{{ID: 1}, {ID: 2}}
	mockPosts.On("GetPostsForModeration", ctx, filter).Return(posts, nil)
	mockReports.On("GetReportReasons", ctx, 1, core.ReportableTypePost).Return([]string{"spam"}, nil)
	mockReports.On("GetReportReasons", ctx, 2, core.ReportableTypePost).Return([]string{"offensive"}, nil)

	svc := moderator.New(nil, mockPosts, mockReports, nil)

	result, err := svc.GetPostsForModeration(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "spam", result[0].Reasons[0])
	assert.Equal(t, "offensive", result[1].Reasons[0])

	mockPosts.AssertExpectations(t)
	mockReports.AssertExpectations(t)
}

func TestGetPostsForModeration_ReportFail_Continues(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockReports := new(mocks.MockReportStore)

	filter := core.FilterASC
	posts := []core.Post{{ID: 1}, {ID: 2}}
	mockPosts.On("GetPostsForModeration", ctx, filter).Return(posts, nil)
	mockReports.On("GetReportReasons", ctx, 1, core.ReportableTypePost).Return(nil, core.ErrGettingReportReasons)
	mockReports.On("GetReportReasons", ctx, 2, core.ReportableTypePost).Return([]string{"spam"}, nil)

	svc := moderator.New(nil, mockPosts, mockReports, nil)

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

	svc := moderator.New(nil, mockPosts, mockReports, nil)

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

	svc := moderator.New(nil, mockPosts, nil, nil)
	err := svc.DeletePost(ctx, 10)

	assert.NoError(t, err)
	mockPosts.AssertExpectations(t)
}

func TestDeletePost_Failure(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockPosts.On("DeletePost", ctx, 99).Return(core.ErrPostNotFound)

	svc := moderator.New(nil, mockPosts, nil, nil)
	err := svc.DeletePost(ctx, 99)

	assert.Error(t, err)
	assert.Equal(t, err, core.ErrPostNotFound)
	mockPosts.AssertExpectations(t)
}

func TestApprovePost_Success(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)
	mockReports.On("DeleteAllReports", ctx, 5, core.ReportableTypePost).Return(nil)
	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(nil)

	svc := moderator.New(nil, mockPosts, mockReports, nil)
	err := svc.ApprovePost(ctx, 5)

	assert.NoError(t, err)
	mockReports.AssertExpectations(t)
}

func TestApprovePost_Failure(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)
	mockReports.On("DeleteAllReports", ctx, 5, core.ReportableTypePost).Return(errors.New("fail"))
	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(nil)

	svc := moderator.New(nil, mockPosts, mockReports, nil)
	err := svc.ApprovePost(ctx, 5)

	assert.Error(t, err)
	mockReports.AssertExpectations(t)
}

func TestApprovePost_Failure_ApprovePostFromModeration(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)

	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(errors.New("approve failed"))

	svc := moderator.New(nil, mockPosts, mockReports, nil)
	err := svc.ApprovePost(ctx, 5)

	assert.Error(t, err)
	assert.EqualError(t, err, "approve failed")
	mockPosts.AssertExpectations(t)
	mockReports.AssertNotCalled(t, "DeleteAllReportsForPost", ctx, 5)
}

func TestBanUser_Success(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{
		UserID:      1,
		ModeratorID: 2,
		ReportID:    func() *int { i := 5; return &i }(),
	}

	activeUser := core.User{
		ID:     1,
		Status: core.Active,
	}

	mockUserStore.On("GetUserByID", ctx, 1).Return(activeUser, nil)
	mockUserStore.On("BanUserWithRecord", ctx, banRecord).Return(nil)

	err := svc.BanUser(ctx, banRecord)

	assert.NoError(t, err)
	mockUserStore.AssertExpectations(t)
}

func TestBanUser_UserNotFound(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{UserID: 999}

	mockUserStore.On("GetUserByID", ctx, 999).Return(core.User{}, core.ErrNoSuchUser)

	err := svc.BanUser(ctx, banRecord)

	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchUser, err)
	mockUserStore.AssertExpectations(t)
	mockUserStore.AssertNotCalled(t, "BanUserWithRecord")
}

func TestBanUser_UserAlreadyBanned(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{UserID: 1}

	bannedUser := core.User{
		ID:     1,
		Status: core.UserBanned,
	}

	mockUserStore.On("GetUserByID", ctx, 1).Return(bannedUser, nil)

	err := svc.BanUser(ctx, banRecord)

	assert.Error(t, err)
	assert.Equal(t, core.ErrUserAlreadyBanned, err)
	mockUserStore.AssertExpectations(t)
	mockUserStore.AssertNotCalled(t, "BanUserWithRecord")
}

func TestBanUser_StoreErrorOnGetUser(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{UserID: 1}

	mockUserStore.On("GetUserByID", ctx, 1).Return(core.User{}, core.ErrNoSuchUser)

	err := svc.BanUser(ctx, banRecord)

	assert.Error(t, err)
	assert.EqualError(t, err, core.ErrNoSuchUser.Error())
	mockUserStore.AssertExpectations(t)
	mockUserStore.AssertNotCalled(t, "BanUserWithRecord")
}

// here i check some spicific database error
func TestBanUser_StoreErrorOnBan(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{UserID: 1}

	activeUser := core.User{
		ID:     1,
		Status: core.Active,
	}

	mockUserStore.On("GetUserByID", ctx, 1).Return(activeUser, nil)
	mockUserStore.On("BanUserWithRecord", ctx, banRecord).Return(errors.New("ban failed"))

	err := svc.BanUser(ctx, banRecord)

	assert.Error(t, err)
	assert.EqualError(t, err, "ban failed")
	mockUserStore.AssertExpectations(t)
}

func TestBanUser_WithNilReportID(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{
		UserID:      1,
		ModeratorID: 2,
		ReportID:    nil,
	}

	activeUser := core.User{
		ID:     1,
		Status: core.Active,
	}

	mockUserStore.On("GetUserByID", ctx, 1).Return(activeUser, nil)
	mockUserStore.On("BanUserWithRecord", ctx, banRecord).Return(nil)

	err := svc.BanUser(ctx, banRecord)

	assert.NoError(t, err)
	mockUserStore.AssertExpectations(t)
}

func TestBanUser_ConcurrentCalls(t *testing.T) {
	ctx := context.TODO()
	mockUserStore := new(mocks.MockUserStore)
	mockModStore := new(mocks.MockModeratorStore)

	svc := moderator.New(mockModStore, nil, nil, mockUserStore)

	banRecord := core.BannedUserRecord{UserID: 1}

	activeUser := core.User{
		ID:     1,
		Status: core.Active,
	}

	bannedUser := core.User{
		ID:     1,
		Status: core.UserBanned,
	}

	mockUserStore.On("GetUserByID", ctx, 1).Return(activeUser, nil).Once()
	mockUserStore.On("BanUserWithRecord", ctx, banRecord).Return(nil).Once()

	mockUserStore.On("GetUserByID", ctx, 1).Return(bannedUser, nil).Once()

	err := svc.BanUser(ctx, banRecord)
	assert.NoError(t, err)

	err = svc.BanUser(ctx, banRecord)
	assert.Error(t, err)
	assert.Equal(t, core.ErrUserAlreadyBanned, err)

	mockUserStore.AssertExpectations(t)
}
