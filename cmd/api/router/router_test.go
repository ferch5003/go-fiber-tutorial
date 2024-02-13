package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestRegister_Successful(t *testing.T) {
	// Given
	app := fiber.New()
	router := NewRouter(app) // Always have the /health endpoint.
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
