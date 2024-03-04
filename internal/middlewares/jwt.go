package middlewares

import (
	"context"
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/session"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

const (
	// fiberJWT returns the fiber contrib JWT middleware to use.
	fiberJWT = "fiber"

	// applicationJWT returns a manual JWT middleware that use Redis.
	applicationJWT = "app"
)

type JWTMiddleware struct {
	ctx            context.Context
	Type           string
	secret         string
	sessionService session.Service
}

func NewJWTMiddleware(ctx context.Context, jwtType, secret string, sessionService session.Service) *JWTMiddleware {
	return &JWTMiddleware{
		ctx:            ctx,
		Type:           jwtType,
		secret:         secret,
		sessionService: sessionService,
	}
}

func (j *JWTMiddleware) GetMiddleware() fiber.Handler {
	switch j.Type {
	case fiberJWT:
		return j.FiberJWTMiddleware(j.secret)
	case applicationJWT:
		return j.ApplicationJWTMiddleware(j.secret)
	default:
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).SendString("server error")
		}
	}
}

func (j *JWTMiddleware) FiberJWTMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

func (j *JWTMiddleware) ApplicationJWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get("Authorization")
		headerToken, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		}

		token, err := jwt.Parse(headerToken, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		}

		err = j.sessionService.SetSession(j.ctx, headerToken, claims)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		}

		return c.Next()
	}
}
