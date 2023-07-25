package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm-crud/database"
	"gorm-crud/models"
	"gorm-crud/utils"
	"gorm.io/gorm"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type LoginResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		response := Response{Status: "Error", Message: "Bad request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		response := Response{Status: "Error", Message: "Username already exists"}
		return c.Status(http.StatusConflict).JSON(response)
	} else if err != gorm.ErrRecordNotFound {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	user.Password = utils.HashPassword(user.Password)
	if err := database.DB.Create(user).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to create user"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := Response{Status: "Success", Message: "User Created"}
	return c.Status(http.StatusOK).JSON(response)
}

func SignIn(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		response := Response{Status: "Error", Message: "Bad request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	var dbUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "User not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	if !utils.VerifyPassword(user.Password, dbUser.Password) {
		response := Response{Status: "Error", Message: "Invalid credentials"}
		return c.Status(http.StatusUnauthorized).JSON(response)
	}
	token, err := utils.CreateToken(user)
	if err != nil {
		response := Response{Status: "Error", Message: "Server error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := LoginResponse{Status: "Success", Message: "Successfully login", Token: token, TokenType: "Bearer"}
	return c.Status(http.StatusOK).JSON(response)
}
