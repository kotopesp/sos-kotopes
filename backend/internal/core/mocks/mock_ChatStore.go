// Code generated by mockery v2.46.2. DO NOT EDIT.

package core

import (
	context "context"

	chat "github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"

	core "github.com/kotopesp/sos-kotopes/internal/core"

	mock "github.com/stretchr/testify/mock"
)

// MockChatStore is an autogenerated mock type for the ChatStore type
type MockChatStore struct {
	mock.Mock
}

type MockChatStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChatStore) EXPECT() *MockChatStore_Expecter {
	return &MockChatStore_Expecter{mock: &_m.Mock}
}

// AddMemberToChat provides a mock function with given fields: ctx, data
func (_m *MockChatStore) AddMemberToChat(ctx context.Context, data core.ChatMember) (core.ChatMember, error) {
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

// MockChatStore_AddMemberToChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddMemberToChat'
type MockChatStore_AddMemberToChat_Call struct {
	*mock.Call
}

// AddMemberToChat is a helper method to define mock.On call
//   - ctx context.Context
//   - data core.ChatMember
func (_e *MockChatStore_Expecter) AddMemberToChat(ctx interface{}, data interface{}) *MockChatStore_AddMemberToChat_Call {
	return &MockChatStore_AddMemberToChat_Call{Call: _e.mock.On("AddMemberToChat", ctx, data)}
}

func (_c *MockChatStore_AddMemberToChat_Call) Run(run func(ctx context.Context, data core.ChatMember)) *MockChatStore_AddMemberToChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(core.ChatMember))
	})
	return _c
}

func (_c *MockChatStore_AddMemberToChat_Call) Return(member core.ChatMember, err error) *MockChatStore_AddMemberToChat_Call {
	_c.Call.Return(member, err)
	return _c
}

func (_c *MockChatStore_AddMemberToChat_Call) RunAndReturn(run func(context.Context, core.ChatMember) (core.ChatMember, error)) *MockChatStore_AddMemberToChat_Call {
	_c.Call.Return(run)
	return _c
}

// CreateChat provides a mock function with given fields: ctx, data
func (_m *MockChatStore) CreateChat(ctx context.Context, data chat.Chat) (chat.Chat, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateChat")
	}

	var r0 chat.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, chat.Chat) (chat.Chat, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, chat.Chat) chat.Chat); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Get(0).(chat.Chat)
	}

	if rf, ok := ret.Get(1).(func(context.Context, chat.Chat) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatStore_CreateChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateChat'
type MockChatStore_CreateChat_Call struct {
	*mock.Call
}

// CreateChat is a helper method to define mock.On call
//   - ctx context.Context
//   - data chat.Chat
func (_e *MockChatStore_Expecter) CreateChat(ctx interface{}, data interface{}) *MockChatStore_CreateChat_Call {
	return &MockChatStore_CreateChat_Call{Call: _e.mock.On("CreateChat", ctx, data)}
}

func (_c *MockChatStore_CreateChat_Call) Run(run func(ctx context.Context, data chat.Chat)) *MockChatStore_CreateChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(chat.Chat))
	})
	return _c
}

func (_c *MockChatStore_CreateChat_Call) Return(_a0 chat.Chat, err error) *MockChatStore_CreateChat_Call {
	_c.Call.Return(_a0, err)
	return _c
}

func (_c *MockChatStore_CreateChat_Call) RunAndReturn(run func(context.Context, chat.Chat) (chat.Chat, error)) *MockChatStore_CreateChat_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteChat provides a mock function with given fields: ctx, chatID, userID
func (_m *MockChatStore) DeleteChat(ctx context.Context, chatID int, userID int) error {
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

// MockChatStore_DeleteChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteChat'
type MockChatStore_DeleteChat_Call struct {
	*mock.Call
}

// DeleteChat is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int
//   - userID int
func (_e *MockChatStore_Expecter) DeleteChat(ctx interface{}, chatID interface{}, userID interface{}) *MockChatStore_DeleteChat_Call {
	return &MockChatStore_DeleteChat_Call{Call: _e.mock.On("DeleteChat", ctx, chatID, userID)}
}

func (_c *MockChatStore_DeleteChat_Call) Run(run func(ctx context.Context, chatID int, userID int)) *MockChatStore_DeleteChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *MockChatStore_DeleteChat_Call) Return(err error) *MockChatStore_DeleteChat_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockChatStore_DeleteChat_Call) RunAndReturn(run func(context.Context, int, int) error) *MockChatStore_DeleteChat_Call {
	_c.Call.Return(run)
	return _c
}

// FindChatByUsers provides a mock function with given fields: ctx, userIds
func (_m *MockChatStore) FindChatByUsers(ctx context.Context, userIds []int) (chat.Chat, error) {
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

// MockChatStore_FindChatByUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindChatByUsers'
type MockChatStore_FindChatByUsers_Call struct {
	*mock.Call
}

// FindChatByUsers is a helper method to define mock.On call
//   - ctx context.Context
//   - userIds []int
func (_e *MockChatStore_Expecter) FindChatByUsers(ctx interface{}, userIds interface{}) *MockChatStore_FindChatByUsers_Call {
	return &MockChatStore_FindChatByUsers_Call{Call: _e.mock.On("FindChatByUsers", ctx, userIds)}
}

func (_c *MockChatStore_FindChatByUsers_Call) Run(run func(ctx context.Context, userIds []int)) *MockChatStore_FindChatByUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int))
	})
	return _c
}

func (_c *MockChatStore_FindChatByUsers_Call) Return(_a0 chat.Chat, _a1 error) *MockChatStore_FindChatByUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockChatStore_FindChatByUsers_Call) RunAndReturn(run func(context.Context, []int) (chat.Chat, error)) *MockChatStore_FindChatByUsers_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllChats provides a mock function with given fields: ctx, sortType, userID
func (_m *MockChatStore) GetAllChats(ctx context.Context, sortType string, userID int) ([]chat.Chat, error) {
	ret := _m.Called(ctx, sortType, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAllChats")
	}

	var r0 []chat.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) ([]chat.Chat, error)); ok {
		return rf(ctx, sortType, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) []chat.Chat); ok {
		r0 = rf(ctx, sortType, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]chat.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, sortType, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockChatStore_GetAllChats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllChats'
type MockChatStore_GetAllChats_Call struct {
	*mock.Call
}

// GetAllChats is a helper method to define mock.On call
//   - ctx context.Context
//   - sortType string
//   - userID int
func (_e *MockChatStore_Expecter) GetAllChats(ctx interface{}, sortType interface{}, userID interface{}) *MockChatStore_GetAllChats_Call {
	return &MockChatStore_GetAllChats_Call{Call: _e.mock.On("GetAllChats", ctx, sortType, userID)}
}

func (_c *MockChatStore_GetAllChats_Call) Run(run func(ctx context.Context, sortType string, userID int)) *MockChatStore_GetAllChats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int))
	})
	return _c
}

func (_c *MockChatStore_GetAllChats_Call) Return(chats []chat.Chat, err error) *MockChatStore_GetAllChats_Call {
	_c.Call.Return(chats, err)
	return _c
}

func (_c *MockChatStore_GetAllChats_Call) RunAndReturn(run func(context.Context, string, int) ([]chat.Chat, error)) *MockChatStore_GetAllChats_Call {
	_c.Call.Return(run)
	return _c
}

// GetChatWithUsersByID provides a mock function with given fields: ctx, chatID, userID
func (_m *MockChatStore) GetChatWithUsersByID(ctx context.Context, chatID int, userID int) (chat.Chat, error) {
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

// MockChatStore_GetChatWithUsersByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChatWithUsersByID'
type MockChatStore_GetChatWithUsersByID_Call struct {
	*mock.Call
}

// GetChatWithUsersByID is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int
//   - userID int
func (_e *MockChatStore_Expecter) GetChatWithUsersByID(ctx interface{}, chatID interface{}, userID interface{}) *MockChatStore_GetChatWithUsersByID_Call {
	return &MockChatStore_GetChatWithUsersByID_Call{Call: _e.mock.On("GetChatWithUsersByID", ctx, chatID, userID)}
}

func (_c *MockChatStore_GetChatWithUsersByID_Call) Run(run func(ctx context.Context, chatID int, userID int)) *MockChatStore_GetChatWithUsersByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *MockChatStore_GetChatWithUsersByID_Call) Return(_a0 chat.Chat, err error) *MockChatStore_GetChatWithUsersByID_Call {
	_c.Call.Return(_a0, err)
	return _c
}

func (_c *MockChatStore_GetChatWithUsersByID_Call) RunAndReturn(run func(context.Context, int, int) (chat.Chat, error)) *MockChatStore_GetChatWithUsersByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockChatStore creates a new instance of MockChatStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockChatStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockChatStore {
	mock := &MockChatStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}