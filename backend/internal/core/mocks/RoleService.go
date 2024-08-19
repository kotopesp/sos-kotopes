// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// RoleService is an autogenerated mock type for the RoleService type
type RoleService struct {
	mock.Mock
}

// DeleteUserRole provides a mock function with given fields: ctx, id, role
func (_m *RoleService) DeleteUserRole(ctx context.Context, id int, role string) error {
	ret := _m.Called(ctx, id, role)

	if len(ret) == 0 {
		panic("no return value specified for deleteUserRole")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, id, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserRoles provides a mock function with given fields: ctx, id
func (_m *RoleService) GetUserRoles(ctx context.Context, id int) ([]core.RoleDetails, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for getUserRoles")
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

// GiveRoleToUser provides a mock function with given fields: ctx, id, role
func (_m *RoleService) GiveRoleToUser(ctx context.Context, id int, role core.GivenRole) (core.RoleDetails, error) {
	ret := _m.Called(ctx, id, role)

	if len(ret) == 0 {
		panic("no return value specified for giveRoleToUser")
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

// UpdateUserRole provides a mock function with given fields: ctx, id, role
func (_m *RoleService) UpdateUserRole(ctx context.Context, id int, role core.UpdateRole) (core.RoleDetails, error) {
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

// NewRoleService creates a new instance of RoleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRoleService(t interface {
	mock.TestingT
	Cleanup(func())
}) *RoleService {
	mock := &RoleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}