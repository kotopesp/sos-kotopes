// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// EntityService is an autogenerated mock type for the EntityService type
type EntityService struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: ctx, params
func (_m *EntityService) GetAll(ctx context.Context, params core.GetAllParams) ([]core.Entity, int, error) {
	ret := _m.Called(ctx, params)

	var r0 []core.Entity
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllParams) ([]core.Entity, int, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.GetAllParams) []core.Entity); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.GetAllParams) int); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, core.GetAllParams) error); ok {
		r2 = rf(ctx, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *EntityService) GetByID(ctx context.Context, id int) (core.Entity, error) {
	ret := _m.Called(ctx, id)

	var r0 core.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (core.Entity, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) core.Entity); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(core.Entity)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEntityService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEntityService creates a new instance of EntityService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEntityService(t mockConstructorTestingTNewEntityService) *EntityService {
	mock := &EntityService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}