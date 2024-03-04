package router

import (
	"context"
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/handler"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/middlewares"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/session"
	"github.com/ferch5003/go-fiber-tutorial/internal/user"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewUserModule = fx.Module("user",
	// Register Repository & Service
	fx.Provide(user.NewRepository),
	fx.Provide(user.NewService),

	// Register Handler
	fx.Provide(handler.NewUserHandler),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewUserRouter,
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type userRouter struct {
	App            fiber.Router
	config         *config.EnvVars
	sessionService session.Service
	Handler        *handler.UserHandler
}

func NewUserRouter(app *fiber.App,
	config *config.EnvVars,
	sessionService session.Service,
	userHandler *handler.UserHandler) Router {
	return &userRouter{
		App:            app,
		config:         config,
		sessionService: sessionService,
		Handler:        userHandler,
	}
}

func (u userRouter) Register() {
	jwtMiddleware := middlewares.NewJWTMiddleware(
		context.Background(),
		u.config.AppSessionType,
		u.config.AppSecretKey,
		u.sessionService,
	)

	u.App.Route("/users", func(api fiber.Router) {
		api.Get("/:id<int>", u.Handler.Get).Name("get")
		api.Post("/register", u.Handler.RegisterUser).Name("register")
		api.Post("/login", u.Handler.LoginUser).Name("login")

		// Using JWT Middleware.
		protectedRoutes := api.Group("", jwtMiddleware.GetMiddleware())
		protectedRoutes.Patch("/:id<int>", u.Handler.Update).Name("update")
		protectedRoutes.Delete("/:id<int>", u.Handler.Delete).Name("delete")
	}, "users.")
}
