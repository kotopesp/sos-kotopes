// Code generated by mockery v2.43.2. DO NOT EDIT.

package core

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// MockKeeperReviewsStore is an autogenerated mock type for the KeeperReviewsStore type
type MockKeeperReviewsStore struct {
	mock.Mock
}

type MockKeeperReviewsStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKeeperReviewsStore) EXPECT() *MockKeeperReviewsStore_Expecter {
	return &MockKeeperReviewsStore_Expecter{mock: &_m.Mock}
}

// CreateReview provides a mock function with given fields: ctx, keeperReview
func (_m *MockKeeperReviewsStore) CreateReview(ctx context.Context, keeperReview core.KeeperReviews) error {
	ret := _m.Called(ctx, keeperReview)

	if len(ret) == 0 {
		panic("no return value specified for CreateReview")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.KeeperReviews) error); ok {
		r0 = rf(ctx, keeperReview)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKeeperReviewsStore_CreateReview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateReview'
type MockKeeperReviewsStore_CreateReview_Call struct {
	*mock.Call
}

// CreateReview is a helper method to define mock.On call
//   - ctx context.Context
//   - keeperReview core.KeeperReviews
func (_e *MockKeeperReviewsStore_Expecter) CreateReview(ctx interface{}, keeperReview interface{}) *MockKeeperReviewsStore_CreateReview_Call {
	return &MockKeeperReviewsStore_CreateReview_Call{Call: _e.mock.On("CreateReview", ctx, keeperReview)}
}

func (_c *MockKeeperReviewsStore_CreateReview_Call) Run(run func(ctx context.Context, keeperReview core.KeeperReviews)) *MockKeeperReviewsStore_CreateReview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.KeeperReviews))
	})
	return _c
}

func (_c *MockKeeperReviewsStore_CreateReview_Call) Return(_a0 error) *MockKeeperReviewsStore_CreateReview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKeeperReviewsStore_CreateReview_Call) RunAndReturn(run func(context.Context, core.KeeperReviews) error) *MockKeeperReviewsStore_CreateReview_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllReviews provides a mock function with given fields: ctx, params
func (_m *MockKeeperReviewsStore) GetAllReviews(ctx context.Context, params core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GetAllReviews")
	}

	var r0 []core.KeeperReviews
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllKeeperReviewsParams) []core.KeeperReviews); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.KeeperReviews)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.GetAllKeeperReviewsParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeeperReviewsStore_GetAllReviews_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllReviews'
type MockKeeperReviewsStore_GetAllReviews_Call struct {
	*mock.Call
}

// GetAllReviews is a helper method to define mock.On call
//   - ctx context.Context
//   - params core.GetAllKeeperReviewsParams
func (_e *MockKeeperReviewsStore_Expecter) GetAllReviews(ctx interface{}, params interface{}) *MockKeeperReviewsStore_GetAllReviews_Call {
	return &MockKeeperReviewsStore_GetAllReviews_Call{Call: _e.mock.On("GetAllReviews", ctx, params)}
}

func (_c *MockKeeperReviewsStore_GetAllReviews_Call) Run(run func(ctx context.Context, params core.GetAllKeeperReviewsParams)) *MockKeeperReviewsStore_GetAllReviews_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.GetAllKeeperReviewsParams))
	})
	return _c
}

func (_c *MockKeeperReviewsStore_GetAllReviews_Call) Return(_a0 []core.KeeperReviews, _a1 error) *MockKeeperReviewsStore_GetAllReviews_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeeperReviewsStore_GetAllReviews_Call) RunAndReturn(run func(context.Context, core.GetAllKeeperReviewsParams) ([]core.KeeperReviews, error)) *MockKeeperReviewsStore_GetAllReviews_Call {
	_c.Call.Return(run)
	return _c
}

// GetByIDReview provides a mock function with given fields: ctx, id
func (_m *MockKeeperReviewsStore) GetByIDReview(ctx context.Context, id int) (core.KeeperReviews, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDReview")
	}

	var r0 core.KeeperReviews
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (core.KeeperReviews, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) core.KeeperReviews); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(core.KeeperReviews)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeeperReviewsStore_GetByIDReview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByIDReview'
type MockKeeperReviewsStore_GetByIDReview_Call struct {
	*mock.Call
}

// GetByIDReview is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockKeeperReviewsStore_Expecter) GetByIDReview(ctx interface{}, id interface{}) *MockKeeperReviewsStore_GetByIDReview_Call {
	return &MockKeeperReviewsStore_GetByIDReview_Call{Call: _e.mock.On("GetByIDReview", ctx, id)}
}

func (_c *MockKeeperReviewsStore_GetByIDReview_Call) Run(run func(ctx context.Context, id int)) *MockKeeperReviewsStore_GetByIDReview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockKeeperReviewsStore_GetByIDReview_Call) Return(_a0 core.KeeperReviews, _a1 error) *MockKeeperReviewsStore_GetByIDReview_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeeperReviewsStore_GetByIDReview_Call) RunAndReturn(run func(context.Context, int) (core.KeeperReviews, error)) *MockKeeperReviewsStore_GetByIDReview_Call {
	_c.Call.Return(run)
	return _c
}

// SoftDeleteReviewByID provides a mock function with given fields: ctx, id
func (_m *MockKeeperReviewsStore) SoftDeleteReviewByID(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for SoftDeleteReviewByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKeeperReviewsStore_SoftDeleteReviewByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SoftDeleteReviewByID'
type MockKeeperReviewsStore_SoftDeleteReviewByID_Call struct {
	*mock.Call
}

// SoftDeleteReviewByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockKeeperReviewsStore_Expecter) SoftDeleteReviewByID(ctx interface{}, id interface{}) *MockKeeperReviewsStore_SoftDeleteReviewByID_Call {
	return &MockKeeperReviewsStore_SoftDeleteReviewByID_Call{Call: _e.mock.On("SoftDeleteReviewByID", ctx, id)}
}

func (_c *MockKeeperReviewsStore_SoftDeleteReviewByID_Call) Run(run func(ctx context.Context, id int)) *MockKeeperReviewsStore_SoftDeleteReviewByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockKeeperReviewsStore_SoftDeleteReviewByID_Call) Return(_a0 error) *MockKeeperReviewsStore_SoftDeleteReviewByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKeeperReviewsStore_SoftDeleteReviewByID_Call) RunAndReturn(run func(context.Context, int) error) *MockKeeperReviewsStore_SoftDeleteReviewByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateReviewByID provides a mock function with given fields: ctx, keeperReview
func (_m *MockKeeperReviewsStore) UpdateReviewByID(ctx context.Context, keeperReview core.UpdateKeeperReviews) (core.KeeperReviews, error) {
	ret := _m.Called(ctx, keeperReview)

	if len(ret) == 0 {
		panic("no return value specified for UpdateReviewByID")
	}

	var r0 core.KeeperReviews
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.UpdateKeeperReviews) (core.KeeperReviews, error)); ok {
		return rf(ctx, keeperReview)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.UpdateKeeperReviews) core.KeeperReviews); ok {
		r0 = rf(ctx, keeperReview)
	} else {
		r0 = ret.Get(0).(core.KeeperReviews)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.UpdateKeeperReviews) error); ok {
		r1 = rf(ctx, keeperReview)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockKeeperReviewsStore_UpdateReviewByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateReviewByID'
type MockKeeperReviewsStore_UpdateReviewByID_Call struct {
	*mock.Call
}

// UpdateReviewByID is a helper method to define mock.On call
//   - ctx context.Context
//   - keeperReview core.UpdateKeeperReviews
func (_e *MockKeeperReviewsStore_Expecter) UpdateReviewByID(ctx interface{}, keeperReview interface{}) *MockKeeperReviewsStore_UpdateReviewByID_Call {
	return &MockKeeperReviewsStore_UpdateReviewByID_Call{Call: _e.mock.On("UpdateReviewByID", ctx, keeperReview)}
}

func (_c *MockKeeperReviewsStore_UpdateReviewByID_Call) Run(run func(ctx context.Context, keeperReview core.UpdateKeeperReviews)) *MockKeeperReviewsStore_UpdateReviewByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.UpdateKeeperReviews))
	})
	return _c
}

func (_c *MockKeeperReviewsStore_UpdateReviewByID_Call) Return(_a0 core.KeeperReviews, _a1 error) *MockKeeperReviewsStore_UpdateReviewByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockKeeperReviewsStore_UpdateReviewByID_Call) RunAndReturn(run func(context.Context, core.UpdateKeeperReviews) (core.KeeperReviews, error)) *MockKeeperReviewsStore_UpdateReviewByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKeeperReviewsStore creates a new instance of MockKeeperReviewsStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKeeperReviewsStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKeeperReviewsStore {
	mock := &MockKeeperReviewsStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
