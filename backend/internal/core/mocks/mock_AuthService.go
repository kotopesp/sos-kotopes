// Code generated by mockery v2.43.2. DO NOT EDIT.

package core

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// MockAuthService is an autogenerated mock type for the AuthService type
type MockAuthService struct {
	mock.Mock
}

type MockAuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthService) EXPECT() *MockAuthService_Expecter {
	return &MockAuthService_Expecter{mock: &_m.Mock}
}

// AuthorizeVK provides a mock function with given fields: ctx, token
func (_m *MockAuthService) AuthorizeVK(ctx context.Context, token string) (*string, *string, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for AuthorizeVK")
	}

	var r0 *string
	var r1 *string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*string, *string, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *string); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *string); ok {
		r1 = rf(ctx, token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, token)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAuthService_AuthorizeVK_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AuthorizeVK'
type MockAuthService_AuthorizeVK_Call struct {
	*mock.Call
}

// AuthorizeVK is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *MockAuthService_Expecter) AuthorizeVK(ctx interface{}, token interface{}) *MockAuthService_AuthorizeVK_Call {
	return &MockAuthService_AuthorizeVK_Call{Call: _e.mock.On("AuthorizeVK", ctx, token)}
}

func (_c *MockAuthService_AuthorizeVK_Call) Run(run func(ctx context.Context, token string)) *MockAuthService_AuthorizeVK_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAuthService_AuthorizeVK_Call) Return(accessToken *string, refreshToken *string, err error) *MockAuthService_AuthorizeVK_Call {
	_c.Call.Return(accessToken, refreshToken, err)
	return _c
}

func (_c *MockAuthService_AuthorizeVK_Call) RunAndReturn(run func(context.Context, string) (*string, *string, error)) *MockAuthService_AuthorizeVK_Call {
	_c.Call.Return(run)
	return _c
}

// ConfigVK provides a mock function with given fields:
func (_m *MockAuthService) ConfigVK() *oauth2.Config {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ConfigVK")
	}

	var r0 *oauth2.Config
	if rf, ok := ret.Get(0).(func() *oauth2.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Config)
		}
	}

	return r0
}

// MockAuthService_ConfigVK_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigVK'
type MockAuthService_ConfigVK_Call struct {
	*mock.Call
}

// ConfigVK is a helper method to define mock.On call
func (_e *MockAuthService_Expecter) ConfigVK() *MockAuthService_ConfigVK_Call {
	return &MockAuthService_ConfigVK_Call{Call: _e.mock.On("ConfigVK")}
}

func (_c *MockAuthService_ConfigVK_Call) Run(run func()) *MockAuthService_ConfigVK_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAuthService_ConfigVK_Call) Return(_a0 *oauth2.Config) *MockAuthService_ConfigVK_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_ConfigVK_Call) RunAndReturn(run func() *oauth2.Config) *MockAuthService_ConfigVK_Call {
	_c.Call.Return(run)
	return _c
}

// GetJWTSecret provides a mock function with given fields:
func (_m *MockAuthService) GetJWTSecret() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetJWTSecret")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockAuthService_GetJWTSecret_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetJWTSecret'
type MockAuthService_GetJWTSecret_Call struct {
	*mock.Call
}

// GetJWTSecret is a helper method to define mock.On call
func (_e *MockAuthService_Expecter) GetJWTSecret() *MockAuthService_GetJWTSecret_Call {
	return &MockAuthService_GetJWTSecret_Call{Call: _e.mock.On("GetJWTSecret")}
}

func (_c *MockAuthService_GetJWTSecret_Call) Run(run func()) *MockAuthService_GetJWTSecret_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAuthService_GetJWTSecret_Call) Return(_a0 []byte) *MockAuthService_GetJWTSecret_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_GetJWTSecret_Call) RunAndReturn(run func() []byte) *MockAuthService_GetJWTSecret_Call {
	_c.Call.Return(run)
	return _c
}

// LoginBasic provides a mock function with given fields: ctx, user
func (_m *MockAuthService) LoginBasic(ctx context.Context, user core.User) (*string, *string, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for LoginBasic")
	}

	var r0 *string
	var r1 *string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, core.User) (*string, *string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.User) *string); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.User) *string); ok {
		r1 = rf(ctx, user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, core.User) error); ok {
		r2 = rf(ctx, user)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAuthService_LoginBasic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoginBasic'
type MockAuthService_LoginBasic_Call struct {
	*mock.Call
}

// LoginBasic is a helper method to define mock.On call
//   - ctx context.Context
//   - user core.User
func (_e *MockAuthService_Expecter) LoginBasic(ctx interface{}, user interface{}) *MockAuthService_LoginBasic_Call {
	return &MockAuthService_LoginBasic_Call{Call: _e.mock.On("LoginBasic", ctx, user)}
}

func (_c *MockAuthService_LoginBasic_Call) Run(run func(ctx context.Context, user core.User)) *MockAuthService_LoginBasic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.User))
	})
	return _c
}

func (_c *MockAuthService_LoginBasic_Call) Return(accessToken *string, refreshToken *string, err error) *MockAuthService_LoginBasic_Call {
	_c.Call.Return(accessToken, refreshToken, err)
	return _c
}

func (_c *MockAuthService_LoginBasic_Call) RunAndReturn(run func(context.Context, core.User) (*string, *string, error)) *MockAuthService_LoginBasic_Call {
	_c.Call.Return(run)
	return _c
}

// Refresh provides a mock function with given fields: ctx, refreshSession
func (_m *MockAuthService) Refresh(ctx context.Context, refreshSession core.RefreshSession) (*string, *string, error) {
	ret := _m.Called(ctx, refreshSession)

	if len(ret) == 0 {
		panic("no return value specified for Refresh")
	}

	var r0 *string
	var r1 *string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, core.RefreshSession) (*string, *string, error)); ok {
		return rf(ctx, refreshSession)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.RefreshSession) *string); ok {
		r0 = rf(ctx, refreshSession)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.RefreshSession) *string); ok {
		r1 = rf(ctx, refreshSession)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*string)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, core.RefreshSession) error); ok {
		r2 = rf(ctx, refreshSession)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockAuthService_Refresh_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Refresh'
type MockAuthService_Refresh_Call struct {
	*mock.Call
}

// Refresh is a helper method to define mock.On call
//   - ctx context.Context
//   - refreshSession core.RefreshSession
func (_e *MockAuthService_Expecter) Refresh(ctx interface{}, refreshSession interface{}) *MockAuthService_Refresh_Call {
	return &MockAuthService_Refresh_Call{Call: _e.mock.On("Refresh", ctx, refreshSession)}
}

func (_c *MockAuthService_Refresh_Call) Run(run func(ctx context.Context, refreshSession core.RefreshSession)) *MockAuthService_Refresh_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.RefreshSession))
	})
	return _c
}

func (_c *MockAuthService_Refresh_Call) Return(accessToken *string, refreshToken *string, err error) *MockAuthService_Refresh_Call {
	_c.Call.Return(accessToken, refreshToken, err)
	return _c
}

func (_c *MockAuthService_Refresh_Call) RunAndReturn(run func(context.Context, core.RefreshSession) (*string, *string, error)) *MockAuthService_Refresh_Call {
	_c.Call.Return(run)
	return _c
}

// SignupBasic provides a mock function with given fields: ctx, user
func (_m *MockAuthService) SignupBasic(ctx context.Context, user core.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for SignupBasic")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, core.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_SignupBasic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignupBasic'
type MockAuthService_SignupBasic_Call struct {
	*mock.Call
}

// SignupBasic is a helper method to define mock.On call
//   - ctx context.Context
//   - user core.User
func (_e *MockAuthService_Expecter) SignupBasic(ctx interface{}, user interface{}) *MockAuthService_SignupBasic_Call {
	return &MockAuthService_SignupBasic_Call{Call: _e.mock.On("SignupBasic", ctx, user)}
}

func (_c *MockAuthService_SignupBasic_Call) Run(run func(ctx context.Context, user core.User)) *MockAuthService_SignupBasic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.User))
	})
	return _c
}

func (_c *MockAuthService_SignupBasic_Call) Return(_a0 error) *MockAuthService_SignupBasic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_SignupBasic_Call) RunAndReturn(run func(context.Context, core.User) error) *MockAuthService_SignupBasic_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuthService creates a new instance of MockAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthService {
	mock := &MockAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
