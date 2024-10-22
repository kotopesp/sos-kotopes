// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	core "github.com/kotopesp/sos-kotopes/internal/core"
	mock "github.com/stretchr/testify/mock"
)

// PostPhotoService is an autogenerated mock type for the PostPhotoService type
type PostPhotoService struct {
	mock.Mock
}

type PostPhotoService_Expecter struct {
	mock *mock.Mock
}

func (_m *PostPhotoService) EXPECT() *PostPhotoService_Expecter {
	return &PostPhotoService_Expecter{mock: &_m.Mock}
}

// AddPhotosPost provides a mock function with given fields: ctx, postID, photos
func (_m *PostPhotoService) AddPhotosPost(ctx context.Context, postID int, photos []core.Photo) ([]core.Photo, error) {
	ret := _m.Called(ctx, postID, photos)

	var r0 []core.Photo
	if rf, ok := ret.Get(0).(func(context.Context, int, []core.Photo) []core.Photo); ok {
		r0 = rf(ctx, postID, photos)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Photo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, []core.Photo) error); ok {
		r1 = rf(ctx, postID, photos)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostPhotoService_AddPhotosPost_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddPhotosPost'
type PostPhotoService_AddPhotosPost_Call struct {
	*mock.Call
}

// AddPhotosPost is a helper method to define mock.On call
//   - ctx context.Context
//   - postID int
//   - photos []core.Photo
func (_e *PostPhotoService_Expecter) AddPhotosPost(ctx interface{}, postID interface{}, photos interface{}) *PostPhotoService_AddPhotosPost_Call {
	return &PostPhotoService_AddPhotosPost_Call{Call: _e.mock.On("AddPhotosPost", ctx, postID, photos)}
}

func (_c *PostPhotoService_AddPhotosPost_Call) Run(run func(ctx context.Context, postID int, photos []core.Photo)) *PostPhotoService_AddPhotosPost_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].([]core.Photo))
	})
	return _c
}

func (_c *PostPhotoService_AddPhotosPost_Call) Return(_a0 []core.Photo, _a1 error) *PostPhotoService_AddPhotosPost_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetPhotosPost provides a mock function with given fields: ctx, postID
func (_m *PostPhotoService) GetPhotosPost(ctx context.Context, postID int) ([]core.Photo, error) {
	ret := _m.Called(ctx, postID)

	var r0 []core.Photo
	if rf, ok := ret.Get(0).(func(context.Context, int) []core.Photo); ok {
		r0 = rf(ctx, postID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]core.Photo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostPhotoService_GetPhotosPost_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPhotosPost'
type PostPhotoService_GetPhotosPost_Call struct {
	*mock.Call
}

// GetPhotosPost is a helper method to define mock.On call
//   - ctx context.Context
//   - postID int
func (_e *PostPhotoService_Expecter) GetPhotosPost(ctx interface{}, postID interface{}) *PostPhotoService_GetPhotosPost_Call {
	return &PostPhotoService_GetPhotosPost_Call{Call: _e.mock.On("GetPhotosPost", ctx, postID)}
}

func (_c *PostPhotoService_GetPhotosPost_Call) Run(run func(ctx context.Context, postID int)) *PostPhotoService_GetPhotosPost_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *PostPhotoService_GetPhotosPost_Call) Return(_a0 []core.Photo, _a1 error) *PostPhotoService_GetPhotosPost_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetPhotosPostByPhotoID provides a mock function with given fields: ctx, postID, photoID
func (_m *PostPhotoService) GetPhotosPostByPhotoID(ctx context.Context, postID int, photoID int) (core.Photo, error) {
	ret := _m.Called(ctx, postID, photoID)

	var r0 core.Photo
	if rf, ok := ret.Get(0).(func(context.Context, int, int) core.Photo); ok {
		r0 = rf(ctx, postID, photoID)
	} else {
		r0 = ret.Get(0).(core.Photo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, postID, photoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PostPhotoService_GetPhotosPostByPhotoID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPhotosPostByPhotoID'
type PostPhotoService_GetPhotosPostByPhotoID_Call struct {
	*mock.Call
}

// GetPhotosPostByPhotoID is a helper method to define mock.On call
//   - ctx context.Context
//   - postID int
//   - photoID int
func (_e *PostPhotoService_Expecter) GetPhotosPostByPhotoID(ctx interface{}, postID interface{}, photoID interface{}) *PostPhotoService_GetPhotosPostByPhotoID_Call {
	return &PostPhotoService_GetPhotosPostByPhotoID_Call{Call: _e.mock.On("GetPhotosPostByPhotoID", ctx, postID, photoID)}
}

func (_c *PostPhotoService_GetPhotosPostByPhotoID_Call) Run(run func(ctx context.Context, postID int, photoID int)) *PostPhotoService_GetPhotosPostByPhotoID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int))
	})
	return _c
}

func (_c *PostPhotoService_GetPhotosPostByPhotoID_Call) Return(_a0 core.Photo, _a1 error) *PostPhotoService_GetPhotosPostByPhotoID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewPostPhotoService interface {
	mock.TestingT
	Cleanup(func())
}

// NewPostPhotoService creates a new instance of PostPhotoService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPostPhotoService(t mockConstructorTestingTNewPostPhotoService) *PostPhotoService {
	mock := &PostPhotoService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
