package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"gorm-crud/utils"
	"net/http"
	"strings"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			unauthorizedStatus := http.StatusUnauthorized
			return c.Status(unauthorizedStatus).JSON(fiber.Map{"Message": "unauthorized", "Status": unauthorizedStatus})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			unauthorizedStatus := http.StatusUnauthorized
			return c.Status(unauthorizedStatus).JSON(fiber.Map{"message": "unauthorized", "status": unauthorizedStatus})
		}
		c.Locals("claims", claims)

		return c.Next()
	}
}
