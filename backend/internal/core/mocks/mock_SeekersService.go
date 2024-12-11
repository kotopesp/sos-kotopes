// Code generated by mockery v2.43.2. DO NOT EDIT.

package core

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// MockSeekersService is an autogenerated mock type for the SeekersService type
type MockSeekersService struct {
	mock.Mock
}

type MockSeekersService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSeekersService) EXPECT() *MockSeekersService_Expecter {
	return &MockSeekersService_Expecter{mock: &_m.Mock}
}

// CreateEquipment provides a mock function with given fields: ctx, equipment
func (_m *MockSeekersService) CreateEquipment(ctx context.Context, equipment core.Equipment) (int, error) {
	ret := _m.Called(ctx, equipment)

	if len(ret) == 0 {
		panic("no return value specified for CreateEquipment")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Equipment) (int, error)); ok {
		return rf(ctx, equipment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.Equipment) int); ok {
		r0 = rf(ctx, equipment)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.Equipment) error); ok {
		r1 = rf(ctx, equipment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSeekersService_CreateEquipment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateEquipment'
type MockSeekersService_CreateEquipment_Call struct {
	*mock.Call
}

// CreateEquipment is a helper method to define mock.On call
//   - ctx context.Context
//   - equipment core.Equipment
func (_e *MockSeekersService_Expecter) CreateEquipment(ctx interface{}, equipment interface{}) *MockSeekersService_CreateEquipment_Call {
	return &MockSeekersService_CreateEquipment_Call{Call: _e.mock.On("CreateEquipment", ctx, equipment)}
}

func (_c *MockSeekersService_CreateEquipment_Call) Run(run func(ctx context.Context, equipment core.Equipment)) *MockSeekersService_CreateEquipment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Equipment))
	})
	return _c
}

func (_c *MockSeekersService_CreateEquipment_Call) Return(_a0 int, _a1 error) *MockSeekersService_CreateEquipment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSeekersService_CreateEquipment_Call) RunAndReturn(run func(context.Context, core.Equipment) (int, error)) *MockSeekersService_CreateEquipment_Call {
	_c.Call.Return(run)
	return _c
}

// CreateSeeker provides a mock function with given fields: ctx, seeker
func (_m *MockSeekersService) CreateSeeker(ctx context.Context, seeker core.Seeker) (core.Seeker, error) {
	ret := _m.Called(ctx, seeker)

	if len(ret) == 0 {
		panic("no return value specified for CreateSeeker")
	}

	var r0 core.Seeker
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.Seeker) (core.Seeker, error)); ok {
		return rf(ctx, seeker)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.Seeker) core.Seeker); ok {
		r0 = rf(ctx, seeker)
	} else {
		r0 = ret.Get(0).(core.Seeker)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.Seeker) error); ok {
		r1 = rf(ctx, seeker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSeekersService_CreateSeeker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSeeker'
type MockSeekersService_CreateSeeker_Call struct {
	*mock.Call
}

// CreateSeeker is a helper method to define mock.On call
//   - ctx context.Context
//   - seeker core.Seeker
func (_e *MockSeekersService_Expecter) CreateSeeker(ctx interface{}, seeker interface{}) *MockSeekersService_CreateSeeker_Call {
	return &MockSeekersService_CreateSeeker_Call{Call: _e.mock.On("CreateSeeker", ctx, seeker)}
}

func (_c *MockSeekersService_CreateSeeker_Call) Run(run func(ctx context.Context, seeker core.Seeker)) *MockSeekersService_CreateSeeker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.Seeker))
	})
	return _c
}

func (_c *MockSeekersService_CreateSeeker_Call) Return(_a0 core.Seeker, _a1 error) *MockSeekersService_CreateSeeker_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSeekersService_CreateSeeker_Call) RunAndReturn(run func(context.Context, core.Seeker) (core.Seeker, error)) *MockSeekersService_CreateSeeker_Call {
	_c.Call.Return(run)
	return _c
}

// GetSeeker provides a mock function with given fields: ctx, id
func (_m *MockSeekersService) GetSeeker(ctx context.Context, id int) (core.Seeker, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetSeeker")
	}

	var r0 core.Seeker
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (core.Seeker, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) core.Seeker); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(core.Seeker)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSeekersService_GetSeeker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSeeker'
type MockSeekersService_GetSeeker_Call struct {
	*mock.Call
}

// GetSeeker is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockSeekersService_Expecter) GetSeeker(ctx interface{}, id interface{}) *MockSeekersService_GetSeeker_Call {
	return &MockSeekersService_GetSeeker_Call{Call: _e.mock.On("GetSeeker", ctx, id)}
}

func (_c *MockSeekersService_GetSeeker_Call) Run(run func(ctx context.Context, id int)) *MockSeekersService_GetSeeker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockSeekersService_GetSeeker_Call) Return(_a0 core.Seeker, _a1 error) *MockSeekersService_GetSeeker_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSeekersService_GetSeeker_Call) RunAndReturn(run func(context.Context, int) (core.Seeker, error)) *MockSeekersService_GetSeeker_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSeeker provides a mock function with given fields: ctx, seeker
func (_m *MockSeekersService) UpdateSeeker(ctx context.Context, seeker core.UpdateSeeker) (core.Seeker, error) {
	ret := _m.Called(ctx, seeker)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSeeker")
	}

	var r0 core.Seeker
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.UpdateSeeker) (core.Seeker, error)); ok {
		return rf(ctx, seeker)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.UpdateSeeker) core.Seeker); ok {
		r0 = rf(ctx, seeker)
	} else {
		r0 = ret.Get(0).(core.Seeker)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.UpdateSeeker) error); ok {
		r1 = rf(ctx, seeker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSeekersService_UpdateSeeker_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSeeker'
type MockSeekersService_UpdateSeeker_Call struct {
	*mock.Call
}

// UpdateSeeker is a helper method to define mock.On call
//   - ctx context.Context
//   - seeker core.UpdateSeeker
func (_e *MockSeekersService_Expecter) UpdateSeeker(ctx interface{}, seeker interface{}) *MockSeekersService_UpdateSeeker_Call {
	return &MockSeekersService_UpdateSeeker_Call{Call: _e.mock.On("UpdateSeeker", ctx, seeker)}
}

func (_c *MockSeekersService_UpdateSeeker_Call) Run(run func(ctx context.Context, seeker core.UpdateSeeker)) *MockSeekersService_UpdateSeeker_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.UpdateSeeker))
	})
	return _c
}

func (_c *MockSeekersService_UpdateSeeker_Call) Return(_a0 core.Seeker, _a1 error) *MockSeekersService_UpdateSeeker_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSeekersService_UpdateSeeker_Call) RunAndReturn(run func(context.Context, core.UpdateSeeker) (core.Seeker, error)) *MockSeekersService_UpdateSeeker_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSeekersService creates a new instance of MockSeekersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSeekersService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSeekersService {
	mock := &MockSeekersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
