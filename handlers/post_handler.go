package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm-crud/database"
	"gorm-crud/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type PostCreateResponse struct {
	Status string                  `json:"status"`
	Post   PostCreateResponseModel `json:"post"`
}
type PostResponseUser struct {
	FullName string `json:"full_name"`
	Gender   string `json:"gender"`
	Username string `json:"username"`
	Verified bool   `json:"verified"`
}
type PostCreateResponseModel struct {
	ID        uint             `json:"id"`
	Content   string           `json:"content"`
	User      PostResponseUser `json:"user"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeleteAt  gorm.DeletedAt   `json:"delete_at"`
}

func PostCreate(c *fiber.Ctx) error {
	post := new(models.Post)
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	if err := c.BodyParser(&post); err != nil {
		response := Response{Status: "Error", Message: "Bad request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	var dbUser models.User
	if err := database.DB.Where("username = ?", username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	post.UserID = int(dbUser.ID)
	if err := database.DB.Create(post).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to create post"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := PostCreateResponse{
		Status: "Success",
		Post: PostCreateResponseModel{
			ID:        post.ID,
			Content:   post.Content,
			User:      getUserById(post.UserID),
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			DeleteAt:  post.DeletedAt,
		},
	}

	return c.Status(http.StatusOK).JSON(response)
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

	var dbUser models.User
	if err := database.DB.Where("username = ?", username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var dbPost models.Post
	if err := database.DB.Where("id = ? AND user_id = ?", id, dbUser.ID).First(&dbPost).Error; err != nil {
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
			User:      getUserById(dbPost.UserID),
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
	var dbUser models.User
	if err := database.DB.Where("username = ?", username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			return c.Status(http.StatusNotFound).JSON(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var dbPost models.Post
	if err := database.DB.Where("id = ? AND user_id = ?", id, dbUser.ID).First(&dbPost).Error; err != nil {
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
			User:      getUserById(dbPost.UserID),
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
			User:      getUserById(dbPost.UserID),
			CreatedAt: dbPost.CreatedAt,
			UpdatedAt: dbPost.UpdatedAt,
			DeleteAt:  dbPost.DeletedAt,
		},
	}

	return c.Status(http.StatusNoContent).JSON(response)
}
func PostGetAll(c *fiber.Ctx) error {
	type PostGetAllModel struct {
		ID           uint             `json:"id"`
		Content      string           `json:"content"`
		User         PostResponseUser `json:"user"`
		CommentCount int              `json:"comment_count"`
		LikeCount    int              `json:"like_count"`
		CreatedAt    time.Time        `json:"created_at"`
	}

	result := database.DB.Exec(`
	        UPDATE posts
	        SET like_count = (
	            SELECT COUNT(*)
	            FROM post_likes
	            WHERE post_likes.post_id = posts.id
	        )
	    `)

	result2 := database.DB.Exec(`
	        UPDATE posts
	        SET comment_count = (
	            SELECT COUNT(*)
	            FROM comments
	            WHERE comments.post_id = posts.id
	        )
	    `)

	if result.Error != nil {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	if result2.Error != nil {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var posts []models.Post
	if err := database.DB.Find(&posts).Error; err != nil {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var responseData []PostGetAllModel
	for _, post := range posts {
		response := PostGetAllModel{
			ID:           post.ID,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt,
			User:         getUserById(post.UserID),
			CommentCount: post.CommentCount,
			LikeCount:    post.LikeCount,
		}
		responseData = append(responseData, response)

	}
	sort.Slice(responseData, func(i, j int) bool {
		return responseData[i].CreatedAt.After(responseData[j].CreatedAt)
	})
	return c.Status(http.StatusOK).JSON(responseData)
}

func getUserById(ID int) PostResponseUser {
	var dbUser models.User
	if err := database.DB.Where("id = ?", ID).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			log.Fatal(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		log.Fatal(response)
	}
	return PostResponseUser{
		FullName: dbUser.FullName,
		Gender:   dbUser.Gender,
		Username: dbUser.Username,
		Verified: dbUser.Verified,
	}
}
func getUserByUsername(username string) PostResponseUser {
	var dbUser models.User
	if err := database.DB.Where("username = ?", username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not found"}
			log.Fatal(response)
		}
		response := Response{Status: "Error", Message: "Database error"}
		log.Fatal(response)
	}
	return PostResponseUser{
		FullName: dbUser.FullName,
		Gender:   dbUser.Gender,
		Username: dbUser.Username,
		Verified: dbUser.Verified,
	}
}

func PostLike(c *fiber.Ctx) error {
	idParam := c.Params("id")
	post_id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	var existingLike models.PostLike
	if err := database.DB.Where("post_id = ? AND username = ?", post_id, username).First(&existingLike).Error; err == nil {
		response := Response{Status: "Error", Message: "Post already Liked"}
		return c.Status(http.StatusConflict).JSON(response)
	} else if err != gorm.ErrRecordNotFound {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	var postLike = models.PostLike{PostID: int(post_id), Username: username}
	if err := database.DB.Create(&postLike).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to like post"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := Response{Status: "Success", Message: "Liked"}
	return c.Status(http.StatusOK).JSON(response)
}

func PostUnlike(c *fiber.Ctx) error {
	idParam := c.Params("id")
	postID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}

	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	var existingLike models.PostLike
	if err := database.DB.Where("post_id = ? AND username = ?", postID, username).First(&existingLike).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := Response{Status: "Error", Message: "Post not Liked"}
			return c.Status(http.StatusNotFound).JSON(response)
		} else {
			response := Response{Status: "Error", Message: "Database error"}
			return c.Status(http.StatusInternalServerError).JSON(response)
		}
	}

	if err := database.DB.Delete(&existingLike).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to unlike post"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	response := Response{Status: "Success", Message: "Unliked"}
	return c.Status(http.StatusOK).JSON(response)
}

func PostComment(c *fiber.Ctx) error {
	comment := new(models.Comment)
	claims := c.Locals("claims").(jwt.MapClaims)
	username := claims["username"].(string)

	if err := c.BodyParser(&comment); err != nil {
		response := Response{Status: "Error", Message: "Bad Request"}
		return c.Status(http.StatusBadRequest).JSON(response)
	}
	comment.Username = username
	if err := database.DB.Create(&comment).Error; err != nil {
		response := Response{Status: "Error", Message: "Failed to like post"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}
	response := Response{Status: "Success", Message: "comment"}
	return c.Status(http.StatusOK).JSON(response)
}

func GetAllCommentById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	post_id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response := Response{Status: "Error", Message: "Post not found"}
		return c.Status(http.StatusNotFound).JSON(response)
	}

	var comments []models.Comment
	if err := database.DB.Where("post_id = ?", post_id).Find(&comments).Error; err != nil {
		response := Response{Status: "Error", Message: "Database error"}
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	type PostGetAllModel struct {
		ID        uint             `json:"id"`
		Comment   string           `json:"Comment"`
		User      PostResponseUser `json:"user"`
		CreatedAt time.Time        `json:"created_at"`
	}
	var responseData []PostGetAllModel
	for _, comment := range comments {
		response := PostGetAllModel{
			ID:        comment.ID,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			User:      getUserByUsername(comment.Username),
		}
		responseData = append(responseData, response)
	}
	return c.Status(http.StatusOK).JSON(responseData)
}
