package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const _todosPath = "/todos"

type todoServiceMock struct {
	mock.Mock
}

func (tsm *todoServiceMock) GetAll(ctx context.Context, userID int) ([]domain.Todo, error) {
	args := tsm.Called(ctx, userID)
	return args.Get(0).([]domain.Todo), args.Error(1)
}

func (tsm *todoServiceMock) Get(ctx context.Context, id int) (domain.Todo, error) {
	args := tsm.Called(ctx, id)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (tsm *todoServiceMock) Save(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	args := tsm.Called(ctx, todo)
	return args.Get(0).(domain.Todo), args.Error(1)
}

func (tsm *todoServiceMock) Completed(ctx context.Context, id int) error {
	args := tsm.Called(ctx, id)
	return args.Error(0)
}

func (tsm *todoServiceMock) Delete(ctx context.Context, id int) error {
	args := tsm.Called(ctx, id)
	return args.Error(0)
}

func createTodoServer(tsm *todoServiceMock) *fiber.App {
	app := fiber.New()

	userHandler := NewTodoHandler(tsm)

	app.Route("/todos", func(api fiber.Router) {
		// Using JWT Middleware.
		protectedRoutes := api.Group("", middlewares.JWTMiddleware(_testSessionConfigs.Secret))
		protectedRoutes.Get("/", userHandler.GetAll).Name("get_all")
	}, "todos.")

	return app
}

func createTodoRequest(method string, url string, isAuthorized bool, body string) (*http.Request, error) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	if isAuthorized {
		token, err := getTestUserSession()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return req, nil
}

func TestTodoHandlerGetAll_Successful(t *testing.T) {
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

	tsm := new(todoServiceMock)
	tsm.On("GetAll", mock.Anything, expectedUserID).Return(expectedTodos, nil)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		_todosPath,
		true,
		"")
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var todos []domain.Todo
	err = json.Unmarshal(body, &todos)
	require.NoError(t, err)

	require.EqualValues(t, expectedTodos, todos)
}

func TestTodoHandlerGetAll_FailsDueToNotAuthenticatedUser(t *testing.T) {
	// Given
	tsm := new(todoServiceMock)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		_todosPath,
		false,
		"")
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	response := string(body)

	require.Contains(t, response, "Invalid or expired JWT")
}

func TestTodoHandlerGetAll_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("sql: no rows in result set")

	tsm := new(todoServiceMock)
	tsm.On("GetAll", mock.Anything, expectedUserID).Return([]domain.Todo{}, expectedErr)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		_todosPath,
		true,
		"")
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Equal(t, expectedErr.Error(), response.Error)
}
