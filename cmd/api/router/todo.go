package router

import (
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/handler"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/middlewares"
	"github.com/ferch5003/go-fiber-tutorial/internal/todo"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewTodoModule = fx.Module("todo",
	// Register Repository & Service
	fx.Provide(todo.NewRepository),
	fx.Provide(todo.NewService),

	// Register Handler
	fx.Provide(handler.NewTodoHandler),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewTodoRouter,
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type todoRouter struct {
	App     fiber.Router
	config  *config.EnvVars
	Handler *handler.TodoHandler
}

func NewTodoRouter(app *fiber.App, config *config.EnvVars, todoHandler *handler.TodoHandler) Router {
	return &todoRouter{
		App:     app,
		config:  config,
		Handler: todoHandler,
	}
}

func (t todoRouter) Register() {
	t.App.Route("/todos", func(api fiber.Router) {
		// Using JWT Middleware.
		protectedRoutes := api.Group("", middlewares.JWTMiddleware(t.config.AppSecretKey))
		protectedRoutes.Get("/", t.Handler.GetAll).Name("get_all")
		protectedRoutes.Get("/:id<int>", t.Handler.Get).Name("get")
		protectedRoutes.Post("/", t.Handler.Save).Name("save")
		protectedRoutes.Patch("/:id<int>/complete", t.Handler.Completed).Name("completed")
		protectedRoutes.Delete("/:id<int>", t.Handler.Delete).Name("delete")
	}, "todos.")
}
