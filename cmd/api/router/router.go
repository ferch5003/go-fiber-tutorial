package router

import (
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Register()
}

type GeneralRouter struct {
	App        fiber.Router
	UserRouter Router
}

func NewRouter(fiber *fiber.App, userRouter Router) *GeneralRouter {
	return &GeneralRouter{
		App:        fiber,
		UserRouter: userRouter,
	}
}

// Register routes.
func (r *GeneralRouter) Register() {
	r.App.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Application is working correctly! 👋")
	})

	// User Routes.
	r.UserRouter.Register()
}
