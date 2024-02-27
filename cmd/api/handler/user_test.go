package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/middlewares"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/jwtauth"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const _usersPath = "/users"

var _testConfigs = &config.EnvVars{
	AppName:      "test",
	AppSecretKey: "test",
}

var _testSessionConfigs *jwtauth.Config

type _jwtInfo struct {
	ID   int
	Name string
}

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

func (usm *userServiceMock) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	args := usm.Called(ctx, email)
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

	userHandler := NewUserHandler(_testConfigs, usm)

	_testSessionConfigs = &jwtauth.Config{
		AppName: _testConfigs.AppName,
		Secret:  _testConfigs.AppSecretKey,
	}

	app.Route("/users", func(api fiber.Router) {
		api.Get("/:id", userHandler.Get).Name("get")
		api.Post("/register", userHandler.RegisterUser).Name("register")
		api.Post("/login", userHandler.LoginUser).Name("login")

		// Using JWT Middleware.
		protectedRoutes := api.Group("", middlewares.JWTMiddleware(userHandler.config.Secret))
		protectedRoutes.Patch("/:id", userHandler.Update).Name("update")
		protectedRoutes.Delete("/:id", userHandler.Delete).Name("delete")
	}, "users.")

	return app
}

func createRequest(method string, url string, userSession *_jwtInfo, body string) (*http.Request, error) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	if userSession != nil {
		token, err := jwtauth.GenerateToken(userSession.ID, userSession.Name, *_testSessionConfigs)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return req, nil
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

	req, err := createRequest(
		fiber.MethodGet,
		fmt.Sprintf("%s/%d", _usersPath, expectedUserID),
		nil,
		"")
	require.NoError(t, err)

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

	req, err := createRequest(fiber.MethodGet, fmt.Sprintf("%s/%s", _usersPath, "is_not_int"), nil, "")
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

func TestUserHandlerGet_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("sql: no rows in result set")

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, expectedUserID).Return(domain.User{}, expectedErr)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodGet, fmt.Sprintf("%s/%s", _usersPath, "1"), nil, "")
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

	req, err := createRequest(fiber.MethodPost, _usersPath+"/register", nil, `{
																	"first_name": "John",
																	"last_name": "Smith",
																	"email": "john@example.com",
																	"password": "12345678"
																}`)
	require.NoError(t, err)

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

	token := showedUser.Token
	showedUser.Token = nil

	require.EqualValues(t, expectedUser, showedUser)
	require.NotNil(t, token)
}

func TestUserHandlerRegisterUser_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPost, _usersPath+"/register", nil, `{invalid_format}`)
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

	req, err := createRequest(fiber.MethodPost, _usersPath+"/register", nil, `{
																	"first_name": "John",
																	"last_name": "Smith",
																	"email": "john@example.com",
																	"password": "12345678"
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

func TestUserHandlerLoginUser_Successful(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	loggedUser := userData
	loggedUser.ID = 1

	expectedUser := showUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
	}

	usm := new(userServiceMock)
	usm.On("GetByEmail", mock.Anything, userData.Email).Return(loggedUser, nil)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPost, _usersPath+"/login", nil, `{
																	"email": "john@example.com",
																	"password": "12345678"
																}`)
	require.NoError(t, err)

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

	token := showedUser.Token
	showedUser.Token = nil

	require.EqualValues(t, expectedUser, showedUser)
	require.NotNil(t, token)
}

func TestUserHandlerLoginUser_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPost, _usersPath+"/login", nil, `{invalid_format}`)
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

func TestUserHandlerLoginUser_FailsDueToServiceError(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	usm := new(userServiceMock)
	usm.On("GetByEmail", mock.Anything, userData.Email).Return(domain.User{}, expectedError)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPost, _usersPath+"/login", nil, `{
																	"email": "john@example.com",
																	"password": "12345678"
																}`)
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

func TestUserHandlerLoginUser_FailsDueToInvalidCredentials(t *testing.T) {
	// Given
	userData := domain.User{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	loggedUser := userData
	loggedUser.ID = 1

	expectedError := errors.New("Email or Password are incorrect.")

	usm := new(userServiceMock)
	usm.On("GetByEmail", mock.Anything, userData.Email).Return(loggedUser, nil)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPost, _usersPath+"/login", nil, `{
																	"email": "john@example.com",
																	"password": "bad_password"
																}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.ErrorContains(t, expectedError, response.Error)
}

func TestUserHandlerUpdate_Successful(t *testing.T) {
	// Given
	userData := domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	authUser := &_jwtInfo{
		ID:   userData.ID,
		Name: fmt.Sprintf("%s %s", userData.FirstName, userData.LastName),
	}

	updatedUser := userData
	updatedUser.LastName = "Second"

	expectedUser := showUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Second",
		Email:     "john@example.com",
	}

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, userData.ID).Return(userData, nil)
	usm.On("Update", mock.Anything, updatedUser).Return(updatedUser, nil)

	server := createServer(usm)

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%d", _usersPath, userData.ID), authUser, `{
																	"last_name": "Second"
																}`)
	require.NoError(t, err)

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

func TestUserHandlerUpdate_FailsDueToInvalidIntParam(t *testing.T) {
	// Given
	usm := new(userServiceMock)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%s", _usersPath, "is_not_int"), authUser, "")
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

func TestUserHandlerUpdate_FailsDueToUnauthorizedUser(t *testing.T) {
	// Given
	expectedErr := errors.New("Updating not user resource")

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, 0).Return(domain.User{}, expectedErr)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%d", _usersPath, 0), authUser, ``)
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Equal(t, expectedErr.Error(), response.Error)
}

func TestUserHandlerUpdate_FailsDueToObtainingUser(t *testing.T) {
	// Given
	expectedErr := errors.New("sql: no rows in result set")

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, 1).Return(domain.User{}, expectedErr)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%d", _usersPath, 1), authUser, ``)
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

func TestUserHandlerUpdate_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	userData := domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, userData.ID).Return(userData, nil)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%d", _usersPath, userData.ID), authUser, `{invalid_format}`)
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

func TestUserHandlerUpdate_FailsDueToServiceError(t *testing.T) {
	// Given
	userData := domain.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@example.com",
		Password:  "12345678",
	}

	updatedUser := userData
	updatedUser.LastName = "Second"

	expectedError := errors.New("Error Code: 1136. Column count doesn't match value count at row 1")

	usm := new(userServiceMock)
	usm.On("Get", mock.Anything, userData.ID).Return(userData, nil)
	usm.On("Update", mock.Anything, updatedUser).Return(domain.User{}, expectedError)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodPatch, fmt.Sprintf("%s/%d", _usersPath, userData.ID), authUser, `{
																	"last_name": "Second"
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

func TestUserHandlerDelete_Successful(t *testing.T) {
	// Given
	expectedUserID := 1

	usm := new(userServiceMock)
	usm.On("Delete", mock.Anything, expectedUserID).Return(nil)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "John Smith",
	}

	req, err := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%d", _usersPath, expectedUserID), authUser, "")
	require.NoError(t, err)

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

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%s", _usersPath, "is_not_int"), authUser, "")
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

func TestUserHandlerDelete_FailsDueToUnauthorizedUser(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("Updating not user resource")

	usm := new(userServiceMock)
	usm.On("Delete", mock.Anything, expectedUserID).Return(expectedErr)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%d", _usersPath, 0), authUser, "")
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Equal(t, expectedErr.Error(), response.Error)
}

func TestUserHandlerDelete_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedUserID := 1
	expectedErr := errors.New("no rows affected")

	usm := new(userServiceMock)
	usm.On("Delete", mock.Anything, expectedUserID).Return(expectedErr)

	server := createServer(usm)

	authUser := &_jwtInfo{
		ID:   1,
		Name: "Failed",
	}

	req, err := createRequest(fiber.MethodDelete, fmt.Sprintf("%s/%d", _usersPath, 1), authUser, "")
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
