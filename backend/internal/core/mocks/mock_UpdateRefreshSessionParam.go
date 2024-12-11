// Code generated by mockery v2.43.2. DO NOT EDIT.

package core

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// MockUpdateRefreshSessionParam is an autogenerated mock type for the UpdateRefreshSessionParam type
type MockUpdateRefreshSessionParam struct {
	mock.Mock
}

type MockUpdateRefreshSessionParam_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdateRefreshSessionParam) EXPECT() *MockUpdateRefreshSessionParam_Expecter {
	return &MockUpdateRefreshSessionParam_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockUpdateRefreshSessionParam) Execute(_a0 *gorm.DB) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUpdateRefreshSessionParam_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockUpdateRefreshSessionParam_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *gorm.DB
func (_e *MockUpdateRefreshSessionParam_Expecter) Execute(_a0 interface{}) *MockUpdateRefreshSessionParam_Execute_Call {
	return &MockUpdateRefreshSessionParam_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockUpdateRefreshSessionParam_Execute_Call) Run(run func(_a0 *gorm.DB)) *MockUpdateRefreshSessionParam_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gorm.DB))
	})
	return _c
}

func (_c *MockUpdateRefreshSessionParam_Execute_Call) Return(_a0 error) *MockUpdateRefreshSessionParam_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUpdateRefreshSessionParam_Execute_Call) RunAndReturn(run func(*gorm.DB) error) *MockUpdateRefreshSessionParam_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUpdateRefreshSessionParam creates a new instance of MockUpdateRefreshSessionParam. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdateRefreshSessionParam(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdateRefreshSessionParam {
	mock := &MockUpdateRefreshSessionParam{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
