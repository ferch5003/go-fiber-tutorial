package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const _usersPath = "/users"

type errorResponse struct {
	Error string `json:"error"`
}

type userServiceMock struct {
	mock.Mock
}

func (usm *userServiceMock) GetAll(ctx context.Context) ([]domain.User, error) {
	args := usm.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (usm *userServiceMock) Get(ctx context.Context, id int) (domain.User, error) {
	args := usm.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (usm *userServiceMock) Save(ctx context.Context, user domain.User) (domain.User, error) {
	args := usm.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (usm *userServiceMock) Update(ctx context.Context, user domain.User) (domain.User, error) {
	args := usm.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (usm *userServiceMock) Delete(ctx context.Context, id int) error {
	args := usm.Called(ctx, id)
	return args.Error(0)
}

func createServer(usm *userServiceMock) *fiber.App {
	app := fiber.New()

	userHandler := NewUserHandler(usm)

	app.Route("/users", func(api fiber.Router) {
		api.Get("/:id", userHandler.Get).Name("get")
		api.Post("/register", userHandler.RegisterUser).Name("register")
		api.Patch("/:id", userHandler.Update).Name("update")
		api.Delete("/:id", userHandler.Delete).Name("delete")
	}, "users.")

	return app
}

func createRequest(method string, url string, body string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req
}

func TestUserHandlerGet_Successful(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	expectedUserID := 1
	expectedUser := showUser{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, expectedUserID).Return(userData, nil)

	server := createServer(usm)

	req := createRequest(fiber.MethodGet, fmt.Sprintf("%s/%d", _usersPath, expectedUserID), "")

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var showedUser showUser
	err = json.Unmarshal(body, &showedUser)
	require.NoError(t, err)

	require.EqualValues(t, expectedUser, showedUser)
}

func TestUserHandlerGet_FailsDueToInvalidIntParam(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	req := createRequest(fiber.MethodGet, fmt.Sprintf("%s/%s", _usersPath, "is_not_int"), "")

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

func TestUserHandlerGet_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("sql: no rows in result set")

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, expectedUserID).Return(domain.User{}, expectedErr)

	server := createServer(usm)

	req := createRequest(fiber.MethodGet, fmt.Sprintf("%s/%s", _usersPath, "1"), "")

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

func TestUserHandlerRegisterUser_Successful(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	savedUser := userData
	savedUser.ID = 1

	expectedUser := showUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	usm := new(userServiceMock)
	usm.On("Save", mock.Anything, userData).Return(savedUser, nil)

	server := createServer(usm)

	req := createRequest(fiber.MethodPost, _usersPath+"/register", `{
																	"first_name": "John",
																	"last_name": "Smith",
																	"email": "john@example.com",
																	"password": "12345678"
																}`)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var showedUser showUser
	err = json.Unmarshal(body, &showedUser)
	require.NoError(t, err)

	require.EqualValues(t, expectedUser, showedUser)
}

func TestUserHandlerRegisterUser_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	req := createRequest(fiber.MethodPost, _usersPath+"/register", `{invalid_format}`)

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

func TestUserHandlerRegisterUser_FailsDueToServiceError(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	usm := new(userServiceMock)
	usm.On("Save", mock.Anything, userData).Return(domain.User{}, expectedError)

	server := createServer(usm)

	req := createRequest(fiber.MethodPost, _usersPath+"/register", `{
																	"first_name": "John",
																	"last_name": "Smith",
																	"email": "john@example.com",
																	"password": "12345678"
																}`)

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

func TestUserHandlerDelete_Successful(t *testing.T) {
	// Given
	expectedUserID := 1

	usm := new(userServiceMock)
	usm.On("Delete", mock.Anything, expectedUserID).Return(nil)

	server := createServer(usm)

	req := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%d", _usersPath, expectedUserID), "")

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	require.NoError(t, err)
}

func TestUserHandlerDelete_FailsDueToInvalidIntParam(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	req := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%s", _usersPath, "is_not_int"), "")

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

func TestUserHandlerDelete_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("no rows affected")

	usm := new(userServiceMock)
	usm.On("Delete", mock.Anything, expectedUserID).Return(expectedErr)

	server := createServer(usm)

	req := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%s", _usersPath, "1"), "")

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
