package seeker

import (
	"context"
	"errors"
	"github.com/kotopesp/sos-kotopes/internal/core"
	mocks "github.com/kotopesp/sos-kotopes/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateSeeker(t *testing.T) {
	t.Parallel()
	mockSeekersStore := mocks.NewMockSeekersStore(t)
	ctx := context.Background()

	seekerService := New(mockSeekersStore)
	testSeeker := core.Seeker{UserID: 1}

	tests := []struct {
		name           string
		setupMocks     func(*mocks.MockSeekersStore)
		inputSeeker    core.Seeker
		expectedSeeker core.Seeker
		expectedErr    error
	}{
		{
			name: "success",
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testSeeker.UserID).
					Return(core.Seeker{}, core.ErrSeekerNotFound).
					Once()
				ms.On("CreateSeeker", ctx, testSeeker).
					Return(testSeeker, nil).
					Once()
			},
			inputSeeker:    testSeeker,
			expectedSeeker: testSeeker,
			expectedErr:    nil,
		},
		{
			name: "seeker already exists",
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testSeeker.UserID).
					Return(testSeeker, nil).
					Once()
			},
			inputSeeker:    testSeeker,
			expectedSeeker: core.Seeker{},
			expectedErr:    core.ErrSeekerExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(mockSeekersStore)

			result, err := seekerService.CreateSeeker(ctx, tt.inputSeeker)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedSeeker, result)
		})
	}
}

func TestService_GetSeeker(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockSeekersStore := mocks.NewMockSeekersStore(t)
	seekerService := New(mockSeekersStore)
	testID := 123
	activeSeeker := core.Seeker{UserID: testID, IsDeleted: false}
	deletedSeeker := core.Seeker{UserID: testID, IsDeleted: true}

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockSeekersStore)
		expectedResult core.Seeker
		expectedError  error
	}{
		{
			name: "success",
			setupMock: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testID).
					Return(activeSeeker, nil).
					Once()
			},
			expectedResult: activeSeeker,
			expectedError:  nil,
		},
		{
			name: "not found error",
			setupMock: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testID).
					Return(core.Seeker{}, errors.New("not found")).
					Once()
			},
			expectedResult: core.Seeker{},
			expectedError:  core.ErrSeekerNotFound,
		},
		{
			name: "deleted seeker",
			setupMock: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testID).
					Return(deletedSeeker, nil).
					Once()
			},
			expectedResult: core.Seeker{},
			expectedError:  core.ErrSeekerDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mockSeekersStore)

			result, err := seekerService.GetSeeker(ctx, testID)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestService_UpdateSeeker(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockSeekersStore := mocks.NewMockSeekersStore(t)
	seekerService := New(mockSeekersStore)
	userID := 123
	baseSeeker := core.Seeker{ID: 1, UserID: userID}

	tests := []struct {
		name           string
		input          core.UpdateSeeker
		setupMocks     func(*mocks.MockSeekersStore)
		expectedResult core.Seeker
		expectedError  error
	}{
		{
			name: "empty update request",
			input: core.UpdateSeeker{
				UserID: &userID,
			},
			setupMocks:    func(ms *mocks.MockSeekersStore) {},
			expectedError: core.ErrEmptyUpdateRequest,
		},
		{
			name: "successful single field update",
			input: core.UpdateSeeker{
				UserID:     &userID,
				AnimalType: ptrString("dog"),
			},
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, userID).
					Return(baseSeeker, nil).
					Once()
				ms.On("UpdateSeeker", ctx, baseSeeker.ID, map[string]interface{}{
					"animal_type": "dog",
				}).Return(baseSeeker, nil).
					Once()
			},
			expectedResult: baseSeeker,
		},
		{
			name: "multiple fields update",
			input: core.UpdateSeeker{
				UserID:   &userID,
				Location: ptrString("Moscow"),
				Price:    ptrInt(1500),
				HaveCar:  ptrBool(true),
			},
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, userID).
					Return(baseSeeker, nil).
					Once()
				ms.On("UpdateSeeker", ctx, baseSeeker.ID, map[string]interface{}{
					"location": "Moscow",
					"price":    1500,
					"have_car": true,
				}).Return(baseSeeker, nil).
					Once()
			},
			expectedResult: baseSeeker,
		},
		{
			name: "user not found",
			input: core.UpdateSeeker{
				UserID:     &userID,
				AnimalType: ptrString("cat"),
			},
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, userID).
					Return(core.Seeker{}, core.ErrSeekerNotFound).
					Once()
			},
			expectedError: core.ErrSeekerNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(mockSeekersStore)

			result, err := seekerService.UpdateSeeker(ctx, tt.input)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrBool(b bool) *bool       { return &b }

func TestService_DeleteSeeker(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockSeekersStore := mocks.NewMockSeekersStore(t)
	seekerService := New(mockSeekersStore)
	testUserID := 123
	testSeeker := core.Seeker{UserID: testUserID}

	tests := []struct {
		name          string
		userID        int
		setupMocks    func(*mocks.MockSeekersStore)
		expectedError error
	}{
		{
			name:   "success",
			userID: testUserID,
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testUserID).
					Return(testSeeker, nil).
					Once()
				ms.On("DeleteSeeker", ctx, testUserID).
					Return(nil).
					Once()
			},
		},
		{
			name:   "seeker not found",
			userID: testUserID,
			setupMocks: func(ms *mocks.MockSeekersStore) {
				ms.On("GetSeeker", ctx, testUserID).
					Return(core.Seeker{}, core.ErrSeekerNotFound).
					Once()
			},
			expectedError: core.ErrSeekerNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(mockSeekersStore)

			err := seekerService.DeleteSeeker(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetAllSeekers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockSeekers := []core.Seeker{
		{UserID: 1},
		{UserID: 2},
	}

	tests := []struct {
		name           string
		inputParams    core.GetAllSeekersParams
		setupMock      func(*mocks.MockSeekersStore)
		expectedResult []core.Seeker
		expectedError  error
	}{
		{
			name: "default sorting parameters",
			inputParams: core.GetAllSeekersParams{
				SortBy:    ptrString(""),
				SortOrder: ptrString(""),
			},
			setupMock: func(ms *mocks.MockSeekersStore) {
				expectedParams := core.GetAllSeekersParams{
					SortBy:    ptrString("created_at"),
					SortOrder: ptrString("desc"),
				}
				ms.On("GetAllSeekers", ctx, expectedParams).
					Return(mockSeekers, nil).
					Once()
			},
			expectedResult: mockSeekers,
		},
		{
			name: "custom sorting parameters",
			inputParams: core.GetAllSeekersParams{
				SortBy:    ptrString("price"),
				SortOrder: ptrString("asc"),
			},
			setupMock: func(ms *mocks.MockSeekersStore) {
				ms.On("GetAllSeekers", ctx, mock.AnythingOfType("core.GetAllSeekersParams")).
					Return(mockSeekers, nil).
					Once()
			},
			expectedResult: mockSeekers,
		},
		{
			name: "error from store",
			inputParams: core.GetAllSeekersParams{
				SortBy:    ptrString("created_at"),
				SortOrder: ptrString("desc"),
			},
			setupMock: func(ms *mocks.MockSeekersStore) {
				ms.On("GetAllSeekers", ctx, mock.Anything).
					Return(nil, errors.New("storage error")).
					Once()
			},
			expectedError: errors.New("storage error"),
		},
		{
			name:        "nil parameters handling",
			inputParams: core.GetAllSeekersParams{},
			setupMock: func(ms *mocks.MockSeekersStore) {
				expectedParams := core.GetAllSeekersParams{
					SortBy:    ptrString("created_at"),
					SortOrder: ptrString("desc"),
				}
				ms.On("GetAllSeekers", ctx, expectedParams).
					Return(mockSeekers, nil).
					Once()
			},
			expectedResult: mockSeekers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := mocks.NewMockSeekersStore(t)
			tt.setupMock(mockStore)
			service := New(mockStore)

			result, err := service.GetAllSeekers(ctx, tt.inputParams)

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockStore.AssertExpectations(t)
		})
	}
}
