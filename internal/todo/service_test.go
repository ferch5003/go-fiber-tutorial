package todo

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

func (mr *mockRepository) GetAll(ctx context.Context, userID int) ([]domain.Todo, error) {
	args := mr.Called(ctx, userID)
	return args.Get(0).([]domain.Todo), args.Error(1)
}

func (mr *mockRepository) Get(ctx context.Context, id int) (domain.Todo, error) {
	args := mr.Called(ctx, id)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (mr *mockRepository) Save(ctx context.Context, todo domain.Todo) (int, error) {
	args := mr.Called(ctx, todo)
	return args.Int(0), args.Error(1)
}

func (mr *mockRepository) Completed(ctx context.Context, id int) error {
	args := mr.Called(ctx, id)
	return args.Error(0)
}

func (mr *mockRepository) Delete(ctx context.Context, id int) error {
	args := mr.Called(ctx, id)
	return args.Error(0)
}

func TestServiceGetAll_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodos := []domain.Todo{
		{
			ID:          1,
			Title:       "Lorem",
			Description: "Ipsum",
			Completed:   false,
			UserID:      expectedUserID,
		},
		{
			ID:          2,
			Title:       "Lorem Ipsum",
			Description: "FLCL",
			Completed:   true,
			UserID:      expectedUserID,
		},
	}

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything, expectedUserID).Return(expectedTodos, nil)

	service := NewService(mr)

	// When
	todos, err := service.GetAll(context.Background(), expectedUserID)

	// Then
	require.NoError(t, err)
	require.Len(t, todos, len(expectedTodos))
	require.EqualValues(t, expectedTodos, todos)
}

func TestServiceGetAll_SuccessfulWithZeroTodos(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodos := make([]domain.Todo, 0)

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything, expectedUserID).Return(expectedTodos, nil)

	service := NewService(mr)

	// When
	todos, err := service.GetAll(context.Background(), expectedUserID)

	// Then
	require.NoError(t, err)
	require.Len(t, todos, len(expectedTodos))
	require.EqualValues(t, expectedTodos, todos)
}

func TestServiceGetAll_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodos := make([]domain.Todo, 0)
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything, expectedUserID).Return(expectedTodos, expectedError)

	service := NewService(mr)

	// When
	todos, err := service.GetAll(context.Background(), expectedUserID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Len(t, todos, 0)
	require.EqualValues(t, expectedTodos, todos)
}

func TestServiceGet_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodo := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      expectedUserID,
	}

	mr := new(mockRepository)
	mr.On("Get", mock.Anything, expectedTodo.ID).Return(expectedTodo, nil)

	service := NewService(mr)

	// When
	todo, err := service.Get(context.Background(), expectedTodo.ID)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedTodo, todo)
}

func TestServiceGet_FailsDueToRepositoryError(t *testing.T) {
	// Given
	nonExistingID := 1
	expectedTodo := domain.Todo{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Get", mock.Anything, nonExistingID).Return(expectedTodo, expectedError)

	service := NewService(mr)

	// When
	todo, err := service.Get(context.Background(), nonExistingID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, expectedTodo, todo)
}

func TestServiceSave_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodo := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      expectedUserID,
	}

	mr := new(mockRepository)
	mr.On("Save", mock.Anything, expectedTodo).Return(expectedTodo.ID, nil)

	service := NewService(mr)

	// When
	todo, err := service.Save(context.Background(), expectedTodo)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedTodo, todo)
}

func TestServiceSave_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedTodo := domain.Todo{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Save", mock.Anything, expectedTodo).Return(0, expectedError)

	service := NewService(mr)

	// When
	todo, err := service.Save(context.Background(), expectedTodo)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, expectedTodo, todo)
}

func TestServiceCompleted_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodo := domain.Todo{
		ID:     1,
		UserID: expectedUserID,
	}

	mr := new(mockRepository)
	mr.On("Completed", mock.Anything, expectedTodo.ID).Return(nil)

	service := NewService(mr)

	// When
	err := service.Completed(context.Background(), expectedTodo.ID)

	// Then
	require.NoError(t, err)
}

func TestServiceCompleted_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodo := domain.Todo{
		ID:     1,
		UserID: expectedUserID,
	}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Completed", mock.Anything, expectedTodo.ID).Return(expectedError)

	service := NewService(mr)

	// When
	err := service.Completed(context.Background(), expectedTodo.ID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
}

func TestServiceDelete_Successful(t *testing.T) {
	// Given
	expectedTodoID := 1

	mr := new(mockRepository)
	mr.On("Delete", mock.Anything, expectedTodoID).Return(nil)

	service := NewService(mr)

	// When
	err := service.Delete(context.Background(), expectedTodoID)

	// Then
	require.NoError(t, err)
}

func TestServiceDelete_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedTodoID := 1
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("Delete", mock.Anything, expectedTodoID).Return(expectedError)

	service := NewService(mr)

	// When
	err := service.Delete(context.Background(), expectedTodoID)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
}
