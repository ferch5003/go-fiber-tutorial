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

	todoHandler := NewTodoHandler(tsm)

	app.Route("/todos", func(api fiber.Router) {
		// Using JWT Middleware.
		protectedRoutes := api.Group("", middlewares.JWTMiddleware(_testSessionConfigs.Secret))
		protectedRoutes.Get("/", todoHandler.GetAll).Name("get_all")
		protectedRoutes.Get("/:id", todoHandler.Get).Name("get")
		protectedRoutes.Post("/", todoHandler.Save).Name("save")
		protectedRoutes.Patch("/:id/complete", todoHandler.Completed).Name("completed")
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

func TestTodoHandlerGet_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedTodo := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      expectedUserID,
	}

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, expectedTodo.ID).Return(expectedTodo, nil)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		fmt.Sprintf("%s/%d", _todosPath, expectedTodo.ID),
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

	var todoData domain.Todo
	err = json.Unmarshal(body, &todoData)
	require.NoError(t, err)

	require.EqualValues(t, expectedTodo, todoData)
}

func TestTodoHandlerGet_FailsDueToInvalidIntParam(t *testing.T) {
	// Given
	tsm := new(todoServiceMock)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		fmt.Sprintf("%s/%s", _todosPath, "is_not_int"),
		true,
		"")
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "failed to convert")
	require.Contains(t, response.Error, "parsing \"is_not_int\"")
	require.Contains(t, response.Error, "invalid syntax")
}

func TestTodoHandlerGet_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedTodoID := 1
	expectedErr := errors.New("sql: no rows in result set")

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, expectedTodoID).Return(domain.Todo{}, expectedErr)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodGet,
		fmt.Sprintf("%s/%s", _todosPath, "1"),
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

func TestTodoHandlerSave_Successful(t *testing.T) {
	// Given
	expectedUserID := 1
	todoData := domain.Todo{
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      expectedUserID,
	}

	expectedTodo := todoData
	expectedTodo.ID = 1

	tsm := new(todoServiceMock)
	tsm.On("Save", mock.Anything, todoData).Return(expectedTodo, nil)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPost,
		_todosPath,
		true,
		`{
					"title": "Lorem",
					"description": "Ipsum"
					}`)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var showedTodo domain.Todo
	err = json.Unmarshal(body, &showedTodo)
	require.NoError(t, err)

	require.EqualValues(t, expectedTodo, showedTodo)
}

func TestTodoHandlerSave_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	tsm := new(todoServiceMock)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPost,
		_todosPath,
		true,
		`{invalid_format}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "invalid character 'i'")
	require.Contains(t, response.Error, "looking for beginning of object key string")
}

func TestTodoHandlerSave_FailsDueToValidations(t *testing.T) {
	// Given
	tsm := new(todoServiceMock)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(fiber.MethodPost,
		_todosPath,
		true,
		`{}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "[Title]: '' | Needs to implement 'required'")
	require.Contains(t, response.Error, "[Description]: '' | Needs to implement 'required'")
}

func TestTodoHandlerSave_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	tsm := new(todoServiceMock)
	tsm.On("Save", mock.Anything, mock.AnythingOfType("domain.Todo")).Return(domain.Todo{}, expectedError)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPost,
		_todosPath,
		true,
		`{
					"title": "Lorem",
					"description": "Ipsum"
					}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "Error Code: 1136")
	require.Contains(t, response.Error, "Column count doesn't match value count at row 1")
}

func TestTodoHandlerCompleted_Successful(t *testing.T) {
	// Given
	todoData := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      1,
	}

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, todoData.ID).Return(todoData, nil)
	tsm.On("Completed", mock.Anything, todoData.ID).Return(nil)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPatch,
		fmt.Sprintf("%s/%d/complete", _todosPath, todoData.ID),
		true,
		`{}`)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response messageResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Message, "Updated successfully")
}

func TestTodoHandlerCompleted_FailsDueToInvalidIntParam(t *testing.T) {
	// Given
	tsm := new(todoServiceMock)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPatch,
		fmt.Sprintf("%s/%s/complete", _todosPath, "is_not_int"),
		true,
		`{}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "failed to convert")
	require.Contains(t, response.Error, "parsing \"is_not_int\"")
	require.Contains(t, response.Error, "invalid syntax")
}

func TestTodoHandlerCompleted_FailsDueToObtainingTodo(t *testing.T) {
	// Given
	expectedTodoID := 1
	expectedErr := errors.New("sql: no rows in result set")

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, expectedTodoID).Return(domain.Todo{}, expectedErr)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPatch,
		fmt.Sprintf("%s/%d/complete", _todosPath, expectedTodoID),
		true,
		`{}`)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Equal(t, expectedErr.Error(), response.Error)
}

func TestTodoHandlerCompleted_FailsDueToUserNotRelatedTodo(t *testing.T) {
	// Given
	todoData := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Scala",
		Completed:   false,
		UserID:      2,
	}

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, todoData.ID).Return(todoData, nil)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPatch,
		fmt.Sprintf("%s/%d/complete", _todosPath, todoData.ID),
		true,
		`{}`)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "This todo is not from this user")
}

func TestTodoHandlerCompleted_FailsDueToServiceError(t *testing.T) {
	// Given
	todoData := domain.Todo{
		ID:          1,
		Title:       "Lorem",
		Description: "Ipsum",
		Completed:   false,
		UserID:      1,
	}

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	tsm := new(todoServiceMock)
	tsm.On("Get", mock.Anything, todoData.ID).Return(todoData, nil)
	tsm.On("Completed", mock.Anything, todoData.ID).Return(expectedError)

	server := createTodoServer(tsm)

	req, err := createTodoRequest(
		fiber.MethodPatch,
		fmt.Sprintf("%s/%d/complete", _todosPath, todoData.ID),
		true,
		`{}`)
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

	require.Contains(t, response.Error, "Error Code: 1136")
	require.Contains(t, response.Error, "Column count doesn't match value count at row 1")
}
