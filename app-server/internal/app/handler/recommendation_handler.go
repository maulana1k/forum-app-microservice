package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maulana1k/forum-app/internal/app/dto"
	"github.com/maulana1k/forum-app/internal/domain/service"
)

type RecommendationHandler struct {
	service service.RecommendationService
}

func NewRecommendationHandler(s service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: s}
}

// GetRecommendedPosts godoc
//
// @Summary      Get recommended posts for user
// @Description  Retrieve personalized posts feed based on userID and topic
// @Tags         Recommendations
// @Produce      json
//
//	@Security		BearerAuth
//
// @Param        user_id query string true  "User ID"
// @Param        topic   query string false "Topic filter"
// @Param        limit   query int    false "Max number of posts"
// @Success      200 {array} dto.PostResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /v1/recommendation/posts [get]
func (h *RecommendationHandler) GetRecommendedPosts(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "user_id is required"})
	}

	topic := c.Query("topic")
	limitStr := c.Query("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	posts, err := h.service.GetRecommendedPosts(userID, topic, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(posts)
}
