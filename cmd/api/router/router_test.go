package router

import (
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

var _testConfigs = &config.EnvVars{
	AppName:      "test",
	AppSecretKey: "test",
}

type mockUserRouter struct {
	mock.Mock
}

func (m *mockUserRouter) Register() {
	m.Called()
}

func TestRegister_Successful(t *testing.T) {
	// Given
	app := fiber.New()

	mur := new(mockUserRouter)
	mur.On("Register")

	router := NewRouter(app, _testConfigs, mur) // Always have the /health endpoint.
	expectedRoute := "/health"
	expectedStatusCode := fiber.StatusOK

	// When
	router.Register() // Register routes.

	req := httptest.NewRequest("GET", expectedRoute, nil)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, resp.StatusCode)
}
