package handler

import (
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/validations"
	"github.com/ferch5003/go-fiber-tutorial/internal/todo"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	validator *validations.XValidator
	service   todo.Service
}

func NewTodoHandler(service todo.Service) *TodoHandler {
	myValidator := validations.NewValidator()

	return &TodoHandler{
		validator: myValidator,
		service:   service,
	}
}

func (h *TodoHandler) GetAll(c *fiber.Ctx) error {
	userID, err := getAuthUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	todos, err := h.service.GetAll(c.Context(), userID)
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

	obtainedTodo, err := h.service.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(obtainedTodo)
}

func (h *TodoHandler) Save(c *fiber.Ctx) error {
	userID, err := getAuthUserID(c)
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

	savedTodo, err := h.service.Save(c.Context(), todoData)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(savedTodo)
}
