// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// CommentService is an autogenerated mock type for the CommentService type
type CommentService struct {
	mock.Mock
}

type CommentService_Expecter struct {
	mock *mock.Mock
}

func (_m *CommentService) EXPECT() *CommentService_Expecter {
	return &CommentService_Expecter{mock: &_m.Mock}
}

// CreateComment provides a mock function with given fields: ctx, comment, userID
func (_m *CommentService) CreateComment(ctx context.Context, comment core.Comment, userID int) (core.Comment, error) {
	ret := _m.Called(ctx, comment, userID)

	var r0 core.Comment
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment, int) core.Comment); ok {
		r0 = rf(ctx, comment, userID)
	} else {
		r0 = ret.Get(0).(core.Comment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, core.Comment, int) error); ok {
		r1 = rf(ctx, comment, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommentService_CreateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateComment'
type CommentService_CreateComment_Call struct {
	*mock.Call
}

// CreateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comment core.Comment
//   - userID int
func (_e *CommentService_Expecter) CreateComment(ctx interface{}, comment interface{}, userID interface{}) *CommentService_CreateComment_Call {
	return &CommentService_CreateComment_Call{Call: _e.mock.On("CreateComment", ctx, comment, userID)}
}

func (_c *CommentService_CreateComment_Call) Run(run func(ctx context.Context, comment core.Comment, userID int)) *CommentService_CreateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment), args[2].(int))
	})
	return _c
}

func (_c *CommentService_CreateComment_Call) Return(data core.Comment, err error) *CommentService_CreateComment_Call {
	_c.Call.Return(data, err)
	return _c
}

// DeleteComment provides a mock function with given fields: ctx, comments
func (_m *CommentService) DeleteComment(ctx context.Context, comments core.Comment) error {
	ret := _m.Called(ctx, comments)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) error); ok {
		r0 = rf(ctx, comments)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CommentService_DeleteComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteComment'
type CommentService_DeleteComment_Call struct {
	*mock.Call
}

// DeleteComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comments core.Comment
func (_e *CommentService_Expecter) DeleteComment(ctx interface{}, comments interface{}) *CommentService_DeleteComment_Call {
	return &CommentService_DeleteComment_Call{Call: _e.mock.On("DeleteComment", ctx, comments)}
}

func (_c *CommentService_DeleteComment_Call) Run(run func(ctx context.Context, comments core.Comment)) *CommentService_DeleteComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment))
	})
	return _c
}

func (_c *CommentService_DeleteComment_Call) Return(_a0 error) *CommentService_DeleteComment_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetAllComments provides a mock function with given fields: ctx, params, userID
func (_m *CommentService) GetAllComments(ctx context.Context, params core.GetAllCommentsParams, userID int) ([]core.Comment, int, error) {
	ret := _m.Called(ctx, params, userID)

	var r0 []core.Comment
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllCommentsParams, int) []core.Comment); ok {
		r0 = rf(ctx, params, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Comment)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, core.GetAllCommentsParams, int) int); ok {
		r1 = rf(ctx, params, userID)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, core.GetAllCommentsParams, int) error); ok {
		r2 = rf(ctx, params, userID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CommentService_GetAllComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllComments'
type CommentService_GetAllComments_Call struct {
	*mock.Call
}

// GetAllComments is a helper method to define mock.On call
//   - ctx context.Context
//   - params core.GetAllCommentsParams
//   - userID int
func (_e *CommentService_Expecter) GetAllComments(ctx interface{}, params interface{}, userID interface{}) *CommentService_GetAllComments_Call {
	return &CommentService_GetAllComments_Call{Call: _e.mock.On("GetAllComments", ctx, params, userID)}
}

func (_c *CommentService_GetAllComments_Call) Run(run func(ctx context.Context, params core.GetAllCommentsParams, userID int)) *CommentService_GetAllComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.GetAllCommentsParams), args[2].(int))
	})
	return _c
}

func (_c *CommentService_GetAllComments_Call) Return(data []core.Comment, total int, err error) *CommentService_GetAllComments_Call {
	_c.Call.Return(data, total, err)
	return _c
}

// UpdateComment provides a mock function with given fields: ctx, comments
func (_m *CommentService) UpdateComment(ctx context.Context, comments core.Comment) (core.Comment, error) {
	ret := _m.Called(ctx, comments)

	var r0 core.Comment
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) core.Comment); ok {
		r0 = rf(ctx, comments)
	} else {
		r0 = ret.Get(0).(core.Comment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, core.Comment) error); ok {
		r1 = rf(ctx, comments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommentService_UpdateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComment'
type CommentService_UpdateComment_Call struct {
	*mock.Call
}

// UpdateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comments core.Comment
func (_e *CommentService_Expecter) UpdateComment(ctx interface{}, comments interface{}) *CommentService_UpdateComment_Call {
	return &CommentService_UpdateComment_Call{Call: _e.mock.On("UpdateComment", ctx, comments)}
}

func (_c *CommentService_UpdateComment_Call) Run(run func(ctx context.Context, comments core.Comment)) *CommentService_UpdateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment))
	})
	return _c
}

func (_c *CommentService_UpdateComment_Call) Return(data core.Comment, err error) *CommentService_UpdateComment_Call {
	_c.Call.Return(data, err)
	return _c
}

type mockConstructorTestingTNewCommentService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCommentService creates a new instance of CommentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCommentService(t mockConstructorTestingTNewCommentService) *CommentService {
	mock := &CommentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
