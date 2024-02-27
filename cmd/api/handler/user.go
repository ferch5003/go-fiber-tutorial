package handler

import (
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/ferch5003/go-fiber-tutorial/internal/domain"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/data"
	"github.com/ferch5003/go-fiber-tutorial/internal/platform/jwtauth"
	"github.com/ferch5003/go-fiber-tutorial/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type UserHandler struct {
	config  *jwtauth.Config
	service user.Service
}

func NewUserHandler(config *config.EnvVars, service user.Service) *UserHandler {
	jwtConfig := &jwtauth.Config{
		AppName: config.AppName,
		Secret:  config.AppSecretKey,
	}

	return &UserHandler{
		config:  jwtConfig,
		service: service,
	}
}

type showUser struct {
	ID        int     `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Token     *string `json:"token,omitempty"`
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	obtainedUser, err := h.service.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var userData showUser
	columns := []string{"ID", "FirstName", "LastName", "Email"}
	data.OverwriteStruct(&userData, obtainedUser, columns)

	return c.Status(fiber.StatusOK).JSON(userData)
}

type registerUser struct {
	FirstName string `json:"first_name" validate:"required,min=5,max=20"`
	LastName  string `json:"last_name" validate:"required,min=5,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var newUser registerUser

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var userData domain.User
	columns := []string{"FirstName", "LastName", "Email", "Password"}
	data.OverwriteStruct(&userData, newUser, columns)

	createdUser, err := h.service.Save(c.Context(), userData)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fullName := fmt.Sprintf("%s %s", createdUser.FirstName, createdUser.LastName)
	token, err := jwtauth.GenerateToken(createdUser.ID, fullName, *h.config)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var showedUser showUser
	columns = []string{"ID", "FirstName", "LastName", "Email"}
	data.OverwriteStruct(&showedUser, createdUser, columns)

	showedUser.Token = &token

	return c.Status(fiber.StatusCreated).JSON(showedUser)
}

type updateUser struct {
	FirstName string `json:"first_name" validate:"min=5,max=20"`
	LastName  string `json:"last_name" validate:"min=5,max=20"`
	Email     string `json:"email" validate:"email"`
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	protectedUser := c.Locals("user").(*jwt.Token)
	claims := protectedUser.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(float64)

	if !ok || id != int(userID) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Updating not user resource",
		})
	}

	obtainedUser, err := h.service.Get(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var userToUpdate updateUser
	if err := c.BodyParser(&userToUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	columns := []string{"FirstName", "LastName", "Email"}
	data.OverwriteStruct(&obtainedUser, userToUpdate, columns)

	updatedUser, err := h.service.Update(c.Context(), obtainedUser)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var showedUser showUser
	columns = []string{"FirstName", "LastName", "Email"}
	data.OverwriteStruct(&showedUser, updatedUser, columns)

	showedUser.ID = id

	return c.Status(fiber.StatusOK).JSON(showedUser)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	protectedUser := c.Locals("user").(*jwt.Token)
	claims := protectedUser.Claims.(jwt.MapClaims)
	userID, ok := claims["sub"].(float64)

	if !ok || id != int(userID) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Updating not user resource",
		})
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
