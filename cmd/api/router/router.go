package router

import "github.com/gofiber/fiber/v2"

type Router struct {
	App fiber.Router
}

func NewRouter(fiber *fiber.App) *Router {
	return &Router{
		App: fiber,
	}
}

// Register routes
func (r *Router) Register() {
	r.App.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Application is working correctly! 👋")
	})
}
