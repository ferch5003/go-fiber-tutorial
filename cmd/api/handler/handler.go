package handler

import (
	"github.com/ferch5003/go-fiber-tutorial/internal/apierrors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getAuthUserID(c *fiber.Ctx) (int, error) {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0, apierrors.ErrAuthUserNotFound
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return 0, apierrors.ErrAuthUserNotFound
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return 0, apierrors.ErrAuthUserNotFound
	}

	return int(userID), nil
}
