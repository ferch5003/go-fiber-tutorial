package router

import (
	"github.com/ferch5003/go-fiber-tutorial/cmd/api/handler"
	"github.com/ferch5003/go-fiber-tutorial/internal/user"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewUserModule = fx.Module("user",
	// Register Repository & Service
	fx.Provide(user.NewRepository),
	fx.Provide(user.NewService),

	// Register Controller
	fx.Provide(handler.NewUserHandler),

	// Register Router
	fx.Provide(NewUserRouter),
)

type userRouter struct {
	App     fiber.Router
	Handler *handler.UserHandler
}

func NewUserRouter(app *fiber.App, userHandler *handler.UserHandler) Router {
	return &userRouter{
		App:     app,
		Handler: userHandler,
	}
}

func (u userRouter) Register() {
	u.App.Route("/users", func(api fiber.Router) {
		api.Get("/:id<int>", u.Handler.Get).Name("get")
		api.Post("/register", u.Handler.RegisterUser).Name("register")
		api.Delete("/:id<int>", u.Handler.Delete).Name("delete")
	}, "users.")
}
