package router

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/handler"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/middlewares"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/session"
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
	App            fiber.Router
	config         *config.EnvVars
	sessionService session.Service
	Handler        *handler.TodoHandler
}

func NewTodoRouter(
	app *fiber.App,
	config *config.EnvVars,
	sessionService session.Service,
	todoHandler *handler.TodoHandler) Router {
	return &todoRouter{
		App:            app,
		config:         config,
		sessionService: sessionService,
		Handler:        todoHandler,
	}
}

func (t todoRouter) Register() {
	jwtMiddleware := middlewares.NewJWTMiddleware(
		context.Background(),
		t.config.AppSessionType,
		t.config.AppSecretKey,
		t.sessionService,
	)

	t.App.Route("/todos", func(api fiber.Router) {
		// Using JWT Middleware.
		protectedRoutes := api.Group("", jwtMiddleware.GetMiddleware())
		protectedRoutes.Get("/", t.Handler.GetAll).Name("get_all")
		protectedRoutes.Get("/:id<int>", t.Handler.Get).Name("get")
		protectedRoutes.Post("/", t.Handler.Save).Name("save")
		protectedRoutes.Patch("/:id<int>/complete", t.Handler.Completed).Name("completed")
		protectedRoutes.Delete("/:id<int>", t.Handler.Delete).Name("delete")
	}, "todos.")
}
