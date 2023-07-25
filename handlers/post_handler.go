package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm-crud/database"
	"gorm-crud/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostCreateResponse struct {
	Status string                  `json:"status"`
	Post   PostCreateResponseModel `json:"post"`
}
type PostCreateResponseModel struct {
	ID        uint           `json:"id"`
	Content   string         `json:"content"`
	Username  string         `json:"username"`
	Tags      []string       `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `json:"delete_at"`
}

func PostCreate(c *fiber.Ctx) error {
	post := new(models.Post)
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	if err := c.BodyParser(&post); err != nil {
		response := Response{Status: "Error", Message: "Bad request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	post.Username = username
	if err := database.DB.Create(post).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to create post"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := PostCreateResponse{
		Status: "Success",
		Post: PostCreateResponseModel{
			ID:        post.ID,
			Content:   post.Content,
			Username:  post.Username,
			Tags:      TagSlice(post.Tags),
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			DeleteAt:  post.DeletedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}
func TagSlice(tags string) []string {
	newTags := strings.Split(strings.ReplaceAll(tags, ",", " "), " ")
	var tagSlice []string
	for _, t := range newTags {
		tagSlice = append(tagSlice, t)
	}
	return tagSlice
}

func PostGet(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	var dbPost models.Post
	if err := database.DB.Where("id = ? AND username = ?", id, username).First(&dbPost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := PostCreateResponse{
		Status: "Success",
		Post: PostCreateResponseModel{
			ID:        dbPost.ID,
			Content:   dbPost.Content,
			Username:  dbPost.Username,
			Tags:      TagSlice(dbPost.Tags),
			CreatedAt: dbPost.CreatedAt,
			UpdatedAt: dbPost.UpdatedAt,
			DeleteAt:  dbPost.DeletedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}

func PostUpdate(c *fiber.Ctx) error {
	type PostUpdateModel struct {
		Content string `json:"content"`
	}
	updatedPost := new(PostUpdateModel)
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	if err := c.BodyParser(&updatedPost); err != nil {
		response := Response{Status: "Error", Message: "Bad request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}

	var dbPost models.Post
	if err := database.DB.Where("id = ? AND username = ?", id, username).First(&dbPost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	if err := database.DB.Model(&dbPost).Updates(models.Post{
		Content: updatedPost.Content,
	}).Error; err != nil {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := PostCreateResponse{
		Status: "Success",
		Post: PostCreateResponseModel{
			ID:        dbPost.ID,
			Content:   dbPost.Content,
			Username:  dbPost.Username,
			Tags:      TagSlice(dbPost.Tags),
			CreatedAt: dbPost.CreatedAt,
			UpdatedAt: dbPost.UpdatedAt,
			DeleteAt:  dbPost.DeletedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
}
func PostDelete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	var dbPost models.Post
	if err := database.DB.Where("id = ? AND username = ?", id, username).First(&dbPost).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	database.DB.Delete(&dbPost, id)

	response := PostCreateResponse{
		Status: "Success",
		Post: PostCreateResponseModel{
			ID:        dbPost.ID,
			Content:   dbPost.Content,
			Username:  dbPost.Username,
			Tags:      TagSlice(dbPost.Tags),
			CreatedAt: dbPost.CreatedAt,
			UpdatedAt: dbPost.UpdatedAt,
			DeleteAt:  dbPost.DeletedAt,
		},
	}

	return c.Status(http.StatusNoContent).JSON(response)
}
