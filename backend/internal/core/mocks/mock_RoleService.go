// Code generated by mockery v2.46.3. DO NOT EDIT.

package core

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// MockRoleService is an autogenerated mock type for the RoleService type
type MockRoleService struct {
	mock.Mock
}

type MockRoleService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRoleService) EXPECT() *MockRoleService_Expecter {
	return &MockRoleService_Expecter{mock: &_m.Mock}
}

// DeleteUserRole provides a mock function with given fields: ctx, id, role
func (_m *MockRoleService) DeleteUserRole(ctx context.Context, id int, role string) error {
	ret := _m.Called(ctx, id, role)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUserRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRoleService_DeleteUserRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUserRole'
type MockRoleService_DeleteUserRole_Call struct {
	*mock.Call
}

// DeleteUserRole is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role string
func (_e *MockRoleService_Expecter) DeleteUserRole(ctx interface{}, id interface{}, role interface{}) *MockRoleService_DeleteUserRole_Call {
	return &MockRoleService_DeleteUserRole_Call{Call: _e.mock.On("DeleteUserRole", ctx, id, role)}
}

func (_c *MockRoleService_DeleteUserRole_Call) Run(run func(ctx context.Context, id int, role string)) *MockRoleService_DeleteUserRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(string))
	})
	return _c
}

func (_c *MockRoleService_DeleteUserRole_Call) Return(err error) *MockRoleService_DeleteUserRole_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockRoleService_DeleteUserRole_Call) RunAndReturn(run func(context.Context, int, string) error) *MockRoleService_DeleteUserRole_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserRoles provides a mock function with given fields: ctx, id
func (_m *MockRoleService) GetUserRoles(ctx context.Context, id int) ([]core.RoleDetails, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUserRoles")
	}

	var r0 []core.RoleDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]core.RoleDetails, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []core.RoleDetails); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.RoleDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRoleService_GetUserRoles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserRoles'
type MockRoleService_GetUserRoles_Call struct {
	*mock.Call
}

// GetUserRoles is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockRoleService_Expecter) GetUserRoles(ctx interface{}, id interface{}) *MockRoleService_GetUserRoles_Call {
	return &MockRoleService_GetUserRoles_Call{Call: _e.mock.On("GetUserRoles", ctx, id)}
}

func (_c *MockRoleService_GetUserRoles_Call) Run(run func(ctx context.Context, id int)) *MockRoleService_GetUserRoles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockRoleService_GetUserRoles_Call) Return(roles []core.RoleDetails, err error) *MockRoleService_GetUserRoles_Call {
	_c.Call.Return(roles, err)
	return _c
}

func (_c *MockRoleService_GetUserRoles_Call) RunAndReturn(run func(context.Context, int) ([]core.RoleDetails, error)) *MockRoleService_GetUserRoles_Call {
	_c.Call.Return(run)
	return _c
}

// GiveRoleToUser provides a mock function with given fields: ctx, id, role
func (_m *MockRoleService) GiveRoleToUser(ctx context.Context, id int, role core.GivenRole) (core.RoleDetails, error) {
	ret := _m.Called(ctx, id, role)

	if len(ret) == 0 {
		panic("no return value specified for GiveRoleToUser")
	}

	var r0 core.RoleDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, core.GivenRole) (core.RoleDetails, error)); ok {
		return rf(ctx, id, role)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, core.GivenRole) core.RoleDetails); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Get(0).(core.RoleDetails)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, core.GivenRole) error); ok {
		r1 = rf(ctx, id, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRoleService_GiveRoleToUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GiveRoleToUser'
type MockRoleService_GiveRoleToUser_Call struct {
	*mock.Call
}

// GiveRoleToUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role core.GivenRole
func (_e *MockRoleService_Expecter) GiveRoleToUser(ctx interface{}, id interface{}, role interface{}) *MockRoleService_GiveRoleToUser_Call {
	return &MockRoleService_GiveRoleToUser_Call{Call: _e.mock.On("GiveRoleToUser", ctx, id, role)}
}

func (_c *MockRoleService_GiveRoleToUser_Call) Run(run func(ctx context.Context, id int, role core.GivenRole)) *MockRoleService_GiveRoleToUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(core.GivenRole))
	})
	return _c
}

func (_c *MockRoleService_GiveRoleToUser_Call) Return(addedRole core.RoleDetails, err error) *MockRoleService_GiveRoleToUser_Call {
	_c.Call.Return(addedRole, err)
	return _c
}

func (_c *MockRoleService_GiveRoleToUser_Call) RunAndReturn(run func(context.Context, int, core.GivenRole) (core.RoleDetails, error)) *MockRoleService_GiveRoleToUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUserRole provides a mock function with given fields: ctx, id, role
func (_m *MockRoleService) UpdateUserRole(ctx context.Context, id int, role core.UpdateRole) (core.RoleDetails, error) {
	ret := _m.Called(ctx, id, role)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserRole")
	}

	var r0 core.RoleDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, core.UpdateRole) (core.RoleDetails, error)); ok {
		return rf(ctx, id, role)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, core.UpdateRole) core.RoleDetails); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Get(0).(core.RoleDetails)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, core.UpdateRole) error); ok {
		r1 = rf(ctx, id, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRoleService_UpdateUserRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUserRole'
type MockRoleService_UpdateUserRole_Call struct {
	*mock.Call
}

// UpdateUserRole is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role core.UpdateRole
func (_e *MockRoleService_Expecter) UpdateUserRole(ctx interface{}, id interface{}, role interface{}) *MockRoleService_UpdateUserRole_Call {
	return &MockRoleService_UpdateUserRole_Call{Call: _e.mock.On("UpdateUserRole", ctx, id, role)}
}

func (_c *MockRoleService_UpdateUserRole_Call) Run(run func(ctx context.Context, id int, role core.UpdateRole)) *MockRoleService_UpdateUserRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(core.UpdateRole))
	})
	return _c
}

func (_c *MockRoleService_UpdateUserRole_Call) Return(updatedRole core.RoleDetails, err error) *MockRoleService_UpdateUserRole_Call {
	_c.Call.Return(updatedRole, err)
	return _c
}

func (_c *MockRoleService_UpdateUserRole_Call) RunAndReturn(run func(context.Context, int, core.UpdateRole) (core.RoleDetails, error)) *MockRoleService_UpdateUserRole_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRoleService creates a new instance of MockRoleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRoleService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRoleService {
	mock := &MockRoleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
