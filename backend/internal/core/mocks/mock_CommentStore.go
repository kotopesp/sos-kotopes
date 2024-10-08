// Code generated by mockery v2.43.2. DO NOT EDIT.

package core

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// MockCommentStore is an autogenerated mock type for the CommentStore type
type MockCommentStore struct {
	mock.Mock
}

type MockCommentStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCommentStore) EXPECT() *MockCommentStore_Expecter {
	return &MockCommentStore_Expecter{mock: &_m.Mock}
}

// CreateComment provides a mock function with given fields: ctx, comment
func (_m *MockCommentStore) CreateComment(ctx context.Context, comment core.Comment) (core.Comment, error) {
	ret := _m.Called(ctx, comment)

	if len(ret) == 0 {
		panic("no return value specified for CreateComment")
	}

	var r0 core.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) (core.Comment, error)); ok {
		return rf(ctx, comment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) core.Comment); ok {
		r0 = rf(ctx, comment)
	} else {
		r0 = ret.Get(0).(core.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.Comment) error); ok {
		r1 = rf(ctx, comment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCommentStore_CreateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateComment'
type MockCommentStore_CreateComment_Call struct {
	*mock.Call
}

// CreateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comment core.Comment
func (_e *MockCommentStore_Expecter) CreateComment(ctx interface{}, comment interface{}) *MockCommentStore_CreateComment_Call {
	return &MockCommentStore_CreateComment_Call{Call: _e.mock.On("CreateComment", ctx, comment)}
}

func (_c *MockCommentStore_CreateComment_Call) Run(run func(ctx context.Context, comment core.Comment)) *MockCommentStore_CreateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment))
	})
	return _c
}

func (_c *MockCommentStore_CreateComment_Call) Return(data core.Comment, err error) *MockCommentStore_CreateComment_Call {
	_c.Call.Return(data, err)
	return _c
}

func (_c *MockCommentStore_CreateComment_Call) RunAndReturn(run func(context.Context, core.Comment) (core.Comment, error)) *MockCommentStore_CreateComment_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteComment provides a mock function with given fields: ctx, comments
func (_m *MockCommentStore) DeleteComment(ctx context.Context, comments core.Comment) error {
	ret := _m.Called(ctx, comments)

	if len(ret) == 0 {
		panic("no return value specified for DeleteComment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) error); ok {
		r0 = rf(ctx, comments)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockCommentStore_DeleteComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteComment'
type MockCommentStore_DeleteComment_Call struct {
	*mock.Call
}

// DeleteComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comments core.Comment
func (_e *MockCommentStore_Expecter) DeleteComment(ctx interface{}, comments interface{}) *MockCommentStore_DeleteComment_Call {
	return &MockCommentStore_DeleteComment_Call{Call: _e.mock.On("DeleteComment", ctx, comments)}
}

func (_c *MockCommentStore_DeleteComment_Call) Run(run func(ctx context.Context, comments core.Comment)) *MockCommentStore_DeleteComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment))
	})
	return _c
}

func (_c *MockCommentStore_DeleteComment_Call) Return(_a0 error) *MockCommentStore_DeleteComment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCommentStore_DeleteComment_Call) RunAndReturn(run func(context.Context, core.Comment) error) *MockCommentStore_DeleteComment_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllComments provides a mock function with given fields: ctx, params
func (_m *MockCommentStore) GetAllComments(ctx context.Context, params core.GetAllCommentsParams) ([]core.Comment, int, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetAllComments")
	}

	var r0 []core.Comment
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllCommentsParams) ([]core.Comment, int, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllCommentsParams) []core.Comment); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.GetAllCommentsParams) int); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, core.GetAllCommentsParams) error); ok {
		r2 = rf(ctx, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockCommentStore_GetAllComments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllComments'
type MockCommentStore_GetAllComments_Call struct {
	*mock.Call
}

// GetAllComments is a helper method to define mock.On call
//   - ctx context.Context
//   - params core.GetAllCommentsParams
func (_e *MockCommentStore_Expecter) GetAllComments(ctx interface{}, params interface{}) *MockCommentStore_GetAllComments_Call {
	return &MockCommentStore_GetAllComments_Call{Call: _e.mock.On("GetAllComments", ctx, params)}
}

func (_c *MockCommentStore_GetAllComments_Call) Run(run func(ctx context.Context, params core.GetAllCommentsParams)) *MockCommentStore_GetAllComments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.GetAllCommentsParams))
	})
	return _c
}

func (_c *MockCommentStore_GetAllComments_Call) Return(data []core.Comment, total int, err error) *MockCommentStore_GetAllComments_Call {
	_c.Call.Return(data, total, err)
	return _c
}

func (_c *MockCommentStore_GetAllComments_Call) RunAndReturn(run func(context.Context, core.GetAllCommentsParams) ([]core.Comment, int, error)) *MockCommentStore_GetAllComments_Call {
	_c.Call.Return(run)
	return _c
}

// GetCommentByID provides a mock function with given fields: ctx, commentID
func (_m *MockCommentStore) GetCommentByID(ctx context.Context, commentID int) (core.Comment, error) {
	ret := _m.Called(ctx, commentID)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentByID")
	}

	var r0 core.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (core.Comment, error)); ok {
		return rf(ctx, commentID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) core.Comment); ok {
		r0 = rf(ctx, commentID)
	} else {
		r0 = ret.Get(0).(core.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, commentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCommentStore_GetCommentByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCommentByID'
type MockCommentStore_GetCommentByID_Call struct {
	*mock.Call
}

// GetCommentByID is a helper method to define mock.On call
//   - ctx context.Context
//   - commentID int
func (_e *MockCommentStore_Expecter) GetCommentByID(ctx interface{}, commentID interface{}) *MockCommentStore_GetCommentByID_Call {
	return &MockCommentStore_GetCommentByID_Call{Call: _e.mock.On("GetCommentByID", ctx, commentID)}
}

func (_c *MockCommentStore_GetCommentByID_Call) Run(run func(ctx context.Context, commentID int)) *MockCommentStore_GetCommentByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockCommentStore_GetCommentByID_Call) Return(data core.Comment, err error) *MockCommentStore_GetCommentByID_Call {
	_c.Call.Return(data, err)
	return _c
}

func (_c *MockCommentStore_GetCommentByID_Call) RunAndReturn(run func(context.Context, int) (core.Comment, error)) *MockCommentStore_GetCommentByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateComment provides a mock function with given fields: ctx, comments
func (_m *MockCommentStore) UpdateComment(ctx context.Context, comments core.Comment) (core.Comment, error) {
	ret := _m.Called(ctx, comments)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComment")
	}

	var r0 core.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) (core.Comment, error)); ok {
		return rf(ctx, comments)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.Comment) core.Comment); ok {
		r0 = rf(ctx, comments)
	} else {
		r0 = ret.Get(0).(core.Comment)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.Comment) error); ok {
		r1 = rf(ctx, comments)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCommentStore_UpdateComment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComment'
type MockCommentStore_UpdateComment_Call struct {
	*mock.Call
}

// UpdateComment is a helper method to define mock.On call
//   - ctx context.Context
//   - comments core.Comment
func (_e *MockCommentStore_Expecter) UpdateComment(ctx interface{}, comments interface{}) *MockCommentStore_UpdateComment_Call {
	return &MockCommentStore_UpdateComment_Call{Call: _e.mock.On("UpdateComment", ctx, comments)}
}

func (_c *MockCommentStore_UpdateComment_Call) Run(run func(ctx context.Context, comments core.Comment)) *MockCommentStore_UpdateComment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Comment))
	})
	return _c
}

func (_c *MockCommentStore_UpdateComment_Call) Return(data core.Comment, err error) *MockCommentStore_UpdateComment_Call {
	_c.Call.Return(data, err)
	return _c
}

func (_c *MockCommentStore_UpdateComment_Call) RunAndReturn(run func(context.Context, core.Comment) (core.Comment, error)) *MockCommentStore_UpdateComment_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCommentStore creates a new instance of MockCommentStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCommentStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCommentStore {
	mock := &MockCommentStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
