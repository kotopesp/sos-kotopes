// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// RoleStore is an autogenerated mock type for the RoleStore type
type RoleStore struct {
	mock.Mock
}

type RoleStore_Expecter struct {
	mock *mock.Mock
}

func (_m *RoleStore) EXPECT() *RoleStore_Expecter {
	return &RoleStore_Expecter{mock: &_m.Mock}
}

// DeleteUserRole provides a mock function with given fields: ctx, id, role
func (_m *RoleStore) DeleteUserRole(ctx context.Context, id int, role string) error {
	ret := _m.Called(ctx, id, role)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleStore_DeleteUserRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUserRole'
type RoleStore_DeleteUserRole_Call struct {
	*mock.Call
}

// DeleteUserRole is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role string
func (_e *RoleStore_Expecter) DeleteUserRole(ctx interface{}, id interface{}, role interface{}) *RoleStore_DeleteUserRole_Call {
	return &RoleStore_DeleteUserRole_Call{Call: _e.mock.On("DeleteUserRole", ctx, id, role)}
}

func (_c *RoleStore_DeleteUserRole_Call) Run(run func(ctx context.Context, id int, role string)) *RoleStore_DeleteUserRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(string))
	})
	return _c
}

func (_c *RoleStore_DeleteUserRole_Call) Return(err error) *RoleStore_DeleteUserRole_Call {
	_c.Call.Return(err)
	return _c
}

// GetUserRoles provides a mock function with given fields: ctx, id
func (_m *RoleStore) GetUserRoles(ctx context.Context, id int) (map[string]core.Role, error) {
	ret := _m.Called(ctx, id)

	var r0 map[string]core.Role
	if rf, ok := ret.Get(0).(func(context.Context, int) map[string]core.Role); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]core.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleStore_GetUserRoles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserRoles'
type RoleStore_GetUserRoles_Call struct {
	*mock.Call
}

// GetUserRoles is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *RoleStore_Expecter) GetUserRoles(ctx interface{}, id interface{}) *RoleStore_GetUserRoles_Call {
	return &RoleStore_GetUserRoles_Call{Call: _e.mock.On("GetUserRoles", ctx, id)}
}

func (_c *RoleStore_GetUserRoles_Call) Run(run func(ctx context.Context, id int)) *RoleStore_GetUserRoles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *RoleStore_GetUserRoles_Call) Return(roles map[string]core.Role, err error) *RoleStore_GetUserRoles_Call {
	_c.Call.Return(roles, err)
	return _c
}

// GiveRoleToUser provides a mock function with given fields: ctx, id, role
func (_m *RoleStore) GiveRoleToUser(ctx context.Context, id int, role core.GivenRole) (core.Role, error) {
	ret := _m.Called(ctx, id, role)

	var r0 core.Role
	if rf, ok := ret.Get(0).(func(context.Context, int, core.GivenRole) core.Role); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Get(0).(core.Role)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, core.GivenRole) error); ok {
		r1 = rf(ctx, id, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleStore_GiveRoleToUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GiveRoleToUser'
type RoleStore_GiveRoleToUser_Call struct {
	*mock.Call
}

// GiveRoleToUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role core.GivenRole
func (_e *RoleStore_Expecter) GiveRoleToUser(ctx interface{}, id interface{}, role interface{}) *RoleStore_GiveRoleToUser_Call {
	return &RoleStore_GiveRoleToUser_Call{Call: _e.mock.On("GiveRoleToUser", ctx, id, role)}
}

func (_c *RoleStore_GiveRoleToUser_Call) Run(run func(ctx context.Context, id int, role core.GivenRole)) *RoleStore_GiveRoleToUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(core.GivenRole))
	})
	return _c
}

func (_c *RoleStore_GiveRoleToUser_Call) Return(addedRole core.Role, err error) *RoleStore_GiveRoleToUser_Call {
	_c.Call.Return(addedRole, err)
	return _c
}

// UpdateUserRole provides a mock function with given fields: ctx, id, role
func (_m *RoleStore) UpdateUserRole(ctx context.Context, id int, role core.UpdateRole) (core.Role, error) {
	ret := _m.Called(ctx, id, role)

	var r0 core.Role
	if rf, ok := ret.Get(0).(func(context.Context, int, core.UpdateRole) core.Role); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Get(0).(core.Role)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, core.UpdateRole) error); ok {
		r1 = rf(ctx, id, role)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleStore_UpdateUserRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUserRole'
type RoleStore_UpdateUserRole_Call struct {
	*mock.Call
}

// UpdateUserRole is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - role core.UpdateRole
func (_e *RoleStore_Expecter) UpdateUserRole(ctx interface{}, id interface{}, role interface{}) *RoleStore_UpdateUserRole_Call {
	return &RoleStore_UpdateUserRole_Call{Call: _e.mock.On("UpdateUserRole", ctx, id, role)}
}

func (_c *RoleStore_UpdateUserRole_Call) Run(run func(ctx context.Context, id int, role core.UpdateRole)) *RoleStore_UpdateUserRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(core.UpdateRole))
	})
	return _c
}

func (_c *RoleStore_UpdateUserRole_Call) Return(updatedRole core.Role, err error) *RoleStore_UpdateUserRole_Call {
	_c.Call.Return(updatedRole, err)
	return _c
}

type mockConstructorTestingTNewRoleStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoleStore creates a new instance of RoleStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoleStore(t mockConstructorTestingTNewRoleStore) *RoleStore {
	mock := &RoleStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
