package router

import (
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/todo"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewTodoModule = fx.Module("todo",
	// Register Repository & Service
	fx.Provide(todo.NewRepository),
	fx.Provide(todo.NewService),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewTodoRouter,
			fx.As(new(Router)),
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type todoRouter struct {
	App    fiber.Router
	config *config.EnvVars
}

func NewTodoRouter(app *fiber.App, config *config.EnvVars) Router {
	return &todoRouter{
		App:    app,
		config: config,
	}
}

func (t todoRouter) Register() {
	t.App.Route("/todos", func(api fiber.Router) {
	}, "todos.")
}
