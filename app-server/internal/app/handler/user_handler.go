package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers godoc
//
//	@Summary		Get all user profiles
//	@Description	Retrieve a list of all user profiles
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	dto.UserResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/users/ [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	// Map GORM models to DTO
	var response dto.UsersResponse
	for _, u := range users {
		response = append(response, dto.UserResponse{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			AvatarURL:   u.AvatarURL,
			Bio:         u.Bio,
			Location:    u.Location,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		})
	}

	return c.JSON(response)
}

// GetUserByID godoc
//
//	@Summary		Get a user by UUID
//	@Description	Retrieve a specific user profile by its UUID
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User UUID"
//	@Success		200	{object}	dto.UserResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "invalid UUID"})
	}

	user, err := h.service.GetUserProfile(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "user not found"})
	}

	response := dto.UserResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Bio:         user.Bio,
		Location:    user.Location,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return c.JSON(response)
}
