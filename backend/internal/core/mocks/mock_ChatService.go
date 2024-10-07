// Code generated by mockery v2.46.2. DO NOT EDIT.

package core

import (
	context "context"

	chat "github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"

	core "github.com/kotopesp/sos-kotopes/internal/core"

	mock "github.com/stretchr/testify/mock"
)

// MockChatService is an autogenerated mock type for the ChatService type
type MockChatService struct {
	mock.Mock
}

type MockChatService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChatService) EXPECT() *MockChatService_Expecter {
	return &MockChatService_Expecter{mock: &_m.Mock}
}

// AddMemberToChat provides a mock function with given fields: ctx, data
func (_m *MockChatService) AddMemberToChat(ctx context.Context, data core.ChatMember) (core.ChatMember, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for AddMemberToChat")
	}

	var r0 core.ChatMember
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, core.ChatMember) (core.ChatMember, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, core.ChatMember) core.ChatMember); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(core.ChatMember)
	}

	if rf, ok := ret.Get(1).(func(context.Context, core.ChatMember) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatService_AddMemberToChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMemberToChat'
type MockChatService_AddMemberToChat_Call struct {
	*mock.Call
}

// AddMemberToChat is a helper method to define mock.On call
//   - ctx context.Context
//   - data core.ChatMember
func (_e *MockChatService_Expecter) AddMemberToChat(ctx interface{}, data interface{}) *MockChatService_AddMemberToChat_Call {
	return &MockChatService_AddMemberToChat_Call{Call: _e.mock.On("AddMemberToChat", ctx, data)}
}

func (_c *MockChatService_AddMemberToChat_Call) Run(run func(ctx context.Context, data core.ChatMember)) *MockChatService_AddMemberToChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.ChatMember))
	})
	return _c
}

func (_c *MockChatService_AddMemberToChat_Call) Return(member core.ChatMember, err error) *MockChatService_AddMemberToChat_Call {
	_c.Call.Return(member, err)
	return _c
}

func (_c *MockChatService_AddMemberToChat_Call) RunAndReturn(run func(context.Context, core.ChatMember) (core.ChatMember, error)) *MockChatService_AddMemberToChat_Call {
	_c.Call.Return(run)
	return _c
}

// CreateChat provides a mock function with given fields: ctx, data, userIds
func (_m *MockChatService) CreateChat(ctx context.Context, data chat.Chat, userIds []int) (chat.Chat, error) {
	ret := _m.Called(ctx, data, userIds)

	if len(ret) == 0 {
		panic("no return value specified for CreateChat")
	}

	var r0 chat.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, chat.Chat, []int) (chat.Chat, error)); ok {
		return rf(ctx, data, userIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, chat.Chat, []int) chat.Chat); ok {
		r0 = rf(ctx, data, userIds)
	} else {
		r0 = ret.Get(0).(chat.Chat)
	}

	if rf, ok := ret.Get(1).(func(context.Context, chat.Chat, []int) error); ok {
		r1 = rf(ctx, data, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatService_CreateChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateChat'
type MockChatService_CreateChat_Call struct {
	*mock.Call
}

// CreateChat is a helper method to define mock.On call
//   - ctx context.Context
//   - data chat.Chat
//   - userIds []int
func (_e *MockChatService_Expecter) CreateChat(ctx interface{}, data interface{}, userIds interface{}) *MockChatService_CreateChat_Call {
	return &MockChatService_CreateChat_Call{Call: _e.mock.On("CreateChat", ctx, data, userIds)}
}

func (_c *MockChatService_CreateChat_Call) Run(run func(ctx context.Context, data chat.Chat, userIds []int)) *MockChatService_CreateChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(chat.Chat), args[2].([]int))
	})
	return _c
}

func (_c *MockChatService_CreateChat_Call) Return(_a0 chat.Chat, err error) *MockChatService_CreateChat_Call {
	_c.Call.Return(_a0, err)
	return _c
}

func (_c *MockChatService_CreateChat_Call) RunAndReturn(run func(context.Context, chat.Chat, []int) (chat.Chat, error)) *MockChatService_CreateChat_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteChat provides a mock function with given fields: ctx, chatID, userID
func (_m *MockChatService) DeleteChat(ctx context.Context, chatID int, userID int) error {
	ret := _m.Called(ctx, chatID, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteChat")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) error); ok {
		r0 = rf(ctx, chatID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockChatService_DeleteChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteChat'
type MockChatService_DeleteChat_Call struct {
	*mock.Call
}

// DeleteChat is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int
//   - userID int
func (_e *MockChatService_Expecter) DeleteChat(ctx interface{}, chatID interface{}, userID interface{}) *MockChatService_DeleteChat_Call {
	return &MockChatService_DeleteChat_Call{Call: _e.mock.On("DeleteChat", ctx, chatID, userID)}
}

func (_c *MockChatService_DeleteChat_Call) Run(run func(ctx context.Context, chatID int, userID int)) *MockChatService_DeleteChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *MockChatService_DeleteChat_Call) Return(err error) *MockChatService_DeleteChat_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockChatService_DeleteChat_Call) RunAndReturn(run func(context.Context, int, int) error) *MockChatService_DeleteChat_Call {
	_c.Call.Return(run)
	return _c
}

// FindChatByUsers provides a mock function with given fields: ctx, userIds
func (_m *MockChatService) FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error) {
	ret := _m.Called(ctx, userIds)

	if len(ret) == 0 {
		panic("no return value specified for FindChatByUsers")
	}

	var r0 chat.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []int) (chat.Chat, error)); ok {
		return rf(ctx, userIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []int) chat.Chat); ok {
		r0 = rf(ctx, userIds)
	} else {
		r0 = ret.Get(0).(chat.Chat)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []int) error); ok {
		r1 = rf(ctx, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatService_FindChatByUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindChatByUsers'
type MockChatService_FindChatByUsers_Call struct {
	*mock.Call
}

// FindChatByUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - userIds []int
func (_e *MockChatService_Expecter) FindChatByUsers(ctx interface{}, userIds interface{}) *MockChatService_FindChatByUsers_Call {
	return &MockChatService_FindChatByUsers_Call{Call: _e.mock.On("FindChatByUsers", ctx, userIds)}
}

func (_c *MockChatService_FindChatByUsers_Call) Run(run func(ctx context.Context, userIds []int)) *MockChatService_FindChatByUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int))
	})
	return _c
}

func (_c *MockChatService_FindChatByUsers_Call) Return(_a0 chat.Chat, _a1 error) *MockChatService_FindChatByUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatService_FindChatByUsers_Call) RunAndReturn(run func(context.Context, []int) (chat.Chat, error)) *MockChatService_FindChatByUsers_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllChats provides a mock function with given fields: ctx, sortType, userID
func (_m *MockChatService) GetAllChats(ctx context.Context, sortType string, userID int) ([]chat.Chat, int, error) {
	ret := _m.Called(ctx, sortType, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllChats")
	}

	var r0 []chat.Chat
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) ([]chat.Chat, int, error)); ok {
		return rf(ctx, sortType, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) []chat.Chat); ok {
		r0 = rf(ctx, sortType, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]chat.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) int); ok {
		r1 = rf(ctx, sortType, userID)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int) error); ok {
		r2 = rf(ctx, sortType, userID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockChatService_GetAllChats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllChats'
type MockChatService_GetAllChats_Call struct {
	*mock.Call
}

// GetAllChats is a helper method to define mock.On call
//   - ctx context.Context
//   - sortType string
//   - userID int
func (_e *MockChatService_Expecter) GetAllChats(ctx interface{}, sortType interface{}, userID interface{}) *MockChatService_GetAllChats_Call {
	return &MockChatService_GetAllChats_Call{Call: _e.mock.On("GetAllChats", ctx, sortType, userID)}
}

func (_c *MockChatService_GetAllChats_Call) Run(run func(ctx context.Context, sortType string, userID int)) *MockChatService_GetAllChats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int))
	})
	return _c
}

func (_c *MockChatService_GetAllChats_Call) Return(chats []chat.Chat, total int, err error) *MockChatService_GetAllChats_Call {
	_c.Call.Return(chats, total, err)
	return _c
}

func (_c *MockChatService_GetAllChats_Call) RunAndReturn(run func(context.Context, string, int) ([]chat.Chat, int, error)) *MockChatService_GetAllChats_Call {
	_c.Call.Return(run)
	return _c
}

// GetChatWithUsersByID provides a mock function with given fields: ctx, chatID, userID
func (_m *MockChatService) GetChatWithUsersByID(ctx context.Context, chatID int, userID int) (chat.Chat, error) {
	ret := _m.Called(ctx, chatID, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetChatWithUsersByID")
	}

	var r0 chat.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) (chat.Chat, error)); ok {
		return rf(ctx, chatID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) chat.Chat); ok {
		r0 = rf(ctx, chatID, userID)
	} else {
		r0 = ret.Get(0).(chat.Chat)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, chatID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatService_GetChatWithUsersByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChatWithUsersByID'
type MockChatService_GetChatWithUsersByID_Call struct {
	*mock.Call
}

// GetChatWithUsersByID is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int
//   - userID int
func (_e *MockChatService_Expecter) GetChatWithUsersByID(ctx interface{}, chatID interface{}, userID interface{}) *MockChatService_GetChatWithUsersByID_Call {
	return &MockChatService_GetChatWithUsersByID_Call{Call: _e.mock.On("GetChatWithUsersByID", ctx, chatID, userID)}
}

func (_c *MockChatService_GetChatWithUsersByID_Call) Run(run func(ctx context.Context, chatID int, userID int)) *MockChatService_GetChatWithUsersByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *MockChatService_GetChatWithUsersByID_Call) Return(_a0 chat.Chat, err error) *MockChatService_GetChatWithUsersByID_Call {
	_c.Call.Return(_a0, err)
	return _c
}

func (_c *MockChatService_GetChatWithUsersByID_Call) RunAndReturn(run func(context.Context, int, int) (chat.Chat, error)) *MockChatService_GetChatWithUsersByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockChatService creates a new instance of MockChatService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockChatService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockChatService {
	mock := &MockChatService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}