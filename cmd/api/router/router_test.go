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

type mockRouter struct {
	mock.Mock
}

func (m *mockRouter) Register() {
	m.Called()
}

func TestRegister_Successful(t *testing.T) {
	// Given
	app := fiber.New()

	mur := new(mockRouter)
	mur.On("Register")

	mtr := new(mockRouter)
	mtr.On("Register")

	router := NewRouter(app, _testConfigs, mur, mtr) // Always have the /health endpoint.
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
