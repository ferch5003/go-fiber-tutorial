package handler

import (
	"github.com/ferch5003/go-fiber-tutorial/internal/apierrors"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/session"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

func getAuthUserID(c *fiber.Ctx, sessionService session.Service, sessionType string) (int, error) {
	switch sessionType {
	case "fiber":
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
	case "app":
		authorizationHeader := c.Get("Authorization")
		headerToken, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			return 0, apierrors.ErrAuthUserNotFound
		}

		data, err := sessionService.GetSession(c.Context(), headerToken)
		if err != nil {
			return 0, apierrors.ErrAuthUserNotFound
		}

		userID, err := strconv.Atoi(data["sub"])
		if err != nil {
			return 0, apierrors.ErrAuthUserNotFound
		}

		return userID, nil
	default:
		return 0, apierrors.ErrAuthUserNotFound
	}
}
