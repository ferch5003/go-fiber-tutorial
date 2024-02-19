package user

import (
	"context"
	"errors"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockRepository struct {
	mock.Mock
}

func (mr *mockRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := mr.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (mr *mockRepository) Get(ctx context.Context, id int) (domain.User, error) {
	args := mr.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (mr *mockRepository) Save(ctx context.Context, user domain.User) (int, error) {
	args := mr.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (mr *mockRepository) Update(ctx context.Context, user domain.User) error {
	args := mr.Called(ctx, user)
	return args.Error(0)
}

func (mr *mockRepository) Delete(ctx context.Context, id int) error {
	args := mr.Called(ctx, id)
	return args.Error(0)
}

func TestServiceGetAll_Successful(t *testing.T) {
	// Given
	expectedUsers := []domain.User{
		{
			ID:        1,
			FirstName: "Jhon",
			LastName:  "Smith",
			Email:     "john@example.com",
		},
		{
			ID:        2,
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
		},
	}

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, users, len(expectedUsers))
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceGetAll_SuccessfulWithZeroUsers(t *testing.T) {
	// Given
	expectedUsers := make([]domain.User, 0)

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, users, len(expectedUsers))
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceGetAll_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUsers := make([]domain.User, 0)
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, expectedError)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Len(t, users, 0)
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceGet_Successful(t *testing.T) {
	// Given
	expectedUser := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	mr := new(mockRepository)
	mr.On("Get", mock.Anything, expectedUser.ID).Return(expectedUser, nil)

	service := NewService(mr)

	// When
	user, err := service.Get(context.Background(), expectedUser.ID)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestServiceGet_FailsDueToRepositoryError(t *testing.T) {
	// Given
	nonExistingID := 1
	expectedUser := domain.User{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Get", mock.Anything, nonExistingID).Return(expectedUser, expectedError)

	service := NewService(mr)

	// When
	user, err := service.Get(context.Background(), nonExistingID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, expectedUser, user)
}

func TestServiceSave_Successful(t *testing.T) {
	// Given
	expectedUser := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345",
	}

	mr := new(mockRepository)
	mr.On("Save", mock.Anything, expectedUser).Return(expectedUser.ID, nil)

	service := NewService(mr)

	// When
	user, err := service.Save(context.Background(), expectedUser)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestServiceSave_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUser := domain.User{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Save", mock.Anything, expectedUser).Return(0, expectedError)

	service := NewService(mr)

	// When
	user, err := service.Save(context.Background(), expectedUser)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, expectedUser, user)
}

func TestServiceUpdate_Successful(t *testing.T) {
	// Given
	expectedUser := domain.User{
		ID:        1,
		FirstName: "Jhon",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	mr := new(mockRepository)
	mr.On("Update", mock.Anything, expectedUser).Return(nil)

	service := NewService(mr)

	// When
	user, err := service.Update(context.Background(), expectedUser)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestServiceUpdate_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUser := domain.User{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Update", mock.Anything, expectedUser).Return(expectedError)

	service := NewService(mr)

	// When
	user, err := service.Update(context.Background(), expectedUser)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, expectedUser, user)
}

func TestServiceDelete_Successful(t *testing.T) {
	// Given
	expectedUserID := 1

	mr := new(mockRepository)
	mr.On("Delete", mock.Anything, expectedUserID).Return(nil)

	service := NewService(mr)

	// When
	err := service.Delete(context.Background(), expectedUserID)

	// Then
	require.NoError(t, err)
}

func TestServiceDelete_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Delete", mock.Anything, expectedUserID).Return(expectedError)

	service := NewService(mr)

	// When
	err := service.Delete(context.Background(), expectedUserID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
}
