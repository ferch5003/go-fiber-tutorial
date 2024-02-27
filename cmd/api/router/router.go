package router

import (
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Register()
}

type GeneralRouter struct {
	App        fiber.Router
	config     *config.EnvVars
	userRouter Router
}

func NewRouter(fiber *fiber.App, config *config.EnvVars, userRouter Router) *GeneralRouter {
	return &GeneralRouter{
		App:        fiber,
		config:     config,
		userRouter: userRouter,
	}
}

// Register routes.
func (r *GeneralRouter) Register() {
	r.App.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Application is working correctly! 👋")
	})

	// User Routes.
	r.userRouter.Register()
}
