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
	svc := moderator.New(mockMod, nil, nil, nil, nil)

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
	svc := moderator.New(mockMod, nil, nil, nil, nil)

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

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)

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

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)

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

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)

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

	svc := moderator.New(nil, mockPosts, nil, nil, nil)
	err := svc.DeletePost(ctx, 10)

	assert.NoError(t, err)
	mockPosts.AssertExpectations(t)
}

func TestDeletePost_Failure(t *testing.T) {
	ctx := context.TODO()
	mockPosts := new(mocks.MockPostStore)
	mockPosts.On("DeletePost", ctx, 99).Return(core.ErrPostNotFound)

	svc := moderator.New(nil, mockPosts, nil, nil, nil)
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

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)
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

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)
	err := svc.ApprovePost(ctx, 5)

	assert.Error(t, err)
	mockReports.AssertExpectations(t)
}

func TestApprovePost_Failure_ApprovePostFromModeration(t *testing.T) {
	ctx := context.TODO()
	mockReports := new(mocks.MockReportStore)
	mockPosts := new(mocks.MockPostStore)

	mockPosts.On("ApprovePostFromModeration", ctx, 5).Return(errors.New("approve failed"))

	svc := moderator.New(nil, mockPosts, mockReports, nil, nil)
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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

	svc := moderator.New(mockModStore, nil, nil, mockUserStore, nil)

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

func TestGetCommentsForModeration_Success(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	filter := core.FilterDESC
	comments := []core.Comment{
		{ID: 1, Content: "comment 1"},
		{ID: 2, Content: "comment 2"},
	}

	mockCommentStore.On("GetCommentsForModeration", ctx, filter).Return(comments, nil)
	mockReportStore.On("GetReportReasons", ctx, 1, core.ReportableTypeComment).Return([]string{"spam"}, nil)
	mockReportStore.On("GetReportReasons", ctx, 2, core.ReportableTypeComment).Return([]string{"offensive"}, nil)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	result, err := svc.GetCommentsForModeration(ctx, filter)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "spam", result[0].Reasons[0])
	assert.Equal(t, "offensive", result[1].Reasons[0])
	assert.Equal(t, comments[0], result[0].Comment)
	assert.Equal(t, comments[1], result[1].Comment)

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
}

func TestGetCommentsForModeration_NoComments(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	filter := core.FilterASC
	mockCommentStore.On("GetCommentsForModeration", ctx, filter).Return([]core.Comment{}, nil)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	result, err := svc.GetCommentsForModeration(ctx, filter)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoCommentsWaitingForModeration, err)
	assert.Nil(t, result)

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertNotCalled(t, "GetReportReasons")
}

func TestGetCommentsForModeration_StoreError(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	filter := core.FilterDESC
	mockCommentStore.On("GetCommentsForModeration", ctx, filter).Return(nil, errors.New("database error"))

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	result, err := svc.GetCommentsForModeration(ctx, filter)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	assert.Nil(t, result)

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertNotCalled(t, "GetReportReasons")
}

func TestGetCommentsForModeration_ReportError_FailsImmediately(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	filter := core.FilterDESC
	comments := []core.Comment{
		{ID: 1, Content: "comment 1"},
		{ID: 2, Content: "comment 2"},
	}

	mockCommentStore.On("GetCommentsForModeration", ctx, filter).Return(comments, nil)
	mockReportStore.On("GetReportReasons", ctx, 1, core.ReportableTypeComment).Return(nil, errors.New("report error"))

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	result, err := svc.GetCommentsForModeration(ctx, filter)
	assert.Error(t, err)
	assert.EqualError(t, err, "report error")
	assert.Nil(t, result)

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
	mockReportStore.AssertNotCalled(t, "GetReportReasons", ctx, 2, core.ReportableTypeComment)
}

func TestDeleteComment_Success(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)

	commentID := 1
	comment := core.Comment{ID: commentID, Content: "test comment"}

	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(comment, nil)
	mockCommentStore.On("DeleteComment", ctx, comment).Return(nil)

	svc := moderator.New(nil, nil, nil, nil, mockCommentStore)

	err := svc.DeleteComment(ctx, commentID)
	assert.NoError(t, err)

	mockCommentStore.AssertExpectations(t)
}

func TestDeleteComment_CommentNotFound(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)

	commentID := 999
	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(core.Comment{}, core.ErrNoSuchComment)

	svc := moderator.New(nil, nil, nil, nil, mockCommentStore)

	err := svc.DeleteComment(ctx, commentID)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchComment, err)

	mockCommentStore.AssertExpectations(t)
	mockCommentStore.AssertNotCalled(t, "DeleteComment")
}

func TestDeleteComment_DeleteError(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)

	commentID := 1
	comment := core.Comment{ID: commentID, Content: "test comment"}

	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(comment, nil)
	mockCommentStore.On("DeleteComment", ctx, comment).Return(errors.New("delete error"))

	svc := moderator.New(nil, nil, nil, nil, mockCommentStore)

	err := svc.DeleteComment(ctx, commentID)
	assert.Error(t, err)
	assert.EqualError(t, err, "delete error")

	mockCommentStore.AssertExpectations(t)
}

func TestApproveComment_Success(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 1
	comment := core.Comment{ID: commentID, Content: "test comment"}

	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(comment, nil)
	mockCommentStore.On("ApproveCommentFromModeration", ctx, commentID).Return(nil)
	mockReportStore.On("DeleteAllReports", ctx, commentID, core.ReportableTypeComment).Return(nil)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.NoError(t, err)

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
}

func TestApproveComment_CommentNotFound(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 999
	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(core.Comment{}, core.ErrNoSuchComment)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchComment, err)

	mockCommentStore.AssertExpectations(t)
	mockCommentStore.AssertNotCalled(t, "ApproveCommentFromModeration")
	mockReportStore.AssertNotCalled(t, "DeleteAllReports")
}

func TestApproveComment_ApproveError(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 1
	comment := core.Comment{ID: commentID, Content: "test comment"}

	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(comment, nil)
	mockCommentStore.On("ApproveCommentFromModeration", ctx, commentID).Return(errors.New("approve error"))

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.Error(t, err)
	assert.EqualError(t, err, "approve error")

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertNotCalled(t, "DeleteAllReports")
}

func TestApproveComment_DeleteReportsError_ButApprovalSucceeds(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 1
	comment := core.Comment{ID: commentID, Content: "test comment"}

	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(comment, nil)
	mockCommentStore.On("ApproveCommentFromModeration", ctx, commentID).Return(nil)
	mockReportStore.On("DeleteAllReports", ctx, commentID, core.ReportableTypeComment).Return(errors.New("delete reports error"))

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.Error(t, err)
	assert.EqualError(t, err, "delete reports error")

	mockCommentStore.AssertExpectations(t)
	mockReportStore.AssertExpectations(t)
}

func TestApproveComment_GetCommentByIDError(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 1
	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(core.Comment{}, core.ErrNoSuchComment)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.Error(t, err)
	assert.EqualError(t, err, core.ErrNoSuchComment.Error())

	mockCommentStore.AssertExpectations(t)
	mockCommentStore.AssertNotCalled(t, "ApproveCommentFromModeration")
	mockReportStore.AssertNotCalled(t, "DeleteAllReports")
}

func TestApproveComment_GetCommentByIDReturnsErrNoSuchComment(t *testing.T) {
	ctx := context.TODO()
	mockCommentStore := new(mocks.MockCommentStore)
	mockReportStore := new(mocks.MockReportStore)

	commentID := 999
	mockCommentStore.On("GetCommentByID", ctx, commentID).Return(core.Comment{}, core.ErrNoSuchComment)

	svc := moderator.New(nil, nil, mockReportStore, nil, mockCommentStore)

	err := svc.ApproveComment(ctx, commentID)
	assert.Error(t, err)
	assert.Equal(t, core.ErrNoSuchComment, err)

	mockCommentStore.AssertExpectations(t)
	mockCommentStore.AssertNotCalled(t, "ApproveCommentFromModeration")
	mockReportStore.AssertNotCalled(t, "DeleteAllReports")
}
