package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserID(c *fiber.Ctx) (uint, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}
	return uint(id), nil
}
