package handler

import (
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/session"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/validations"
	"github.com/ferch5003/go-fiber-tutorial/internal/todo"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	validator      *validations.XValidator
	sessionType    string
	todoService    todo.Service
	sessionService session.Service
}

func NewTodoHandler(cfg *config.EnvVars, todoService todo.Service, sessionService session.Service) *TodoHandler {
	myValidator := validations.NewValidator()

	return &TodoHandler{
		validator:      myValidator,
		sessionType:    cfg.AppSessionType,
		todoService:    todoService,
		sessionService: sessionService,
	}
}

func (h *TodoHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getAuthUserID(c, h.sessionService, h.sessionType)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	todos, err := h.todoService.GetAll(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(todos)
}

func (h *TodoHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	obtainedTodo, err := h.todoService.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(obtainedTodo)
}

func (h *TodoHandler) Save(c *fiber.Ctx) error {
	userID, err := getAuthUserID(c, h.sessionService, h.sessionType)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var todoData domain.Todo
	if err := c.BodyParser(&todoData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Relate todo to authenticated user.
	todoData.UserID = userID

	todoValidations := h.validator.GetValidations(todoData)
	if todoValidations != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": todoValidations,
		})
	}

	savedTodo, err := h.todoService.Save(c.Context(), todoData)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(savedTodo)
}

func (h *TodoHandler) Completed(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	obtainedTodo, err := h.todoService.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, err := getAuthUserID(c, h.sessionService, h.sessionType)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if obtainedTodo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "This todo is not from this user",
		})
	}

	if err := h.todoService.Completed(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Updated successfully",
	})
}

func (h *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	obtainedTodo, err := h.todoService.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, err := getAuthUserID(c, h.sessionService, h.sessionType)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if obtainedTodo.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "This todo is not from this user",
		})
	}

	if err := h.todoService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo deleted successfully",
	})
}
