package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm-crud/database"
	"gorm-crud/handlers"
	"gorm-crud/middlewares"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{AppName: "Gorm Crud Post App"})
	app.Use(logger.New())
	app.Use(cors.New())
	database.ConnectDB()
	api := app.Group("/api")

	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)

	auth := middlewares.AuthMiddleware()
	api.Use(auth)

	api.Get("/users", handlers.GetAllUser)

	api.Get("/posts", handlers.PostGetAll)

	api.Post("/post_like/:id", handlers.PostLike)
	api.Post("/post_unlike/:id", handlers.PostUnlike)

	api.Post("/post/comment", handlers.PostComment)
	api.Get("/post/comments/:id", handlers.GetAllCommentById)

	api.Post("/post", handlers.PostCreate)       //  Create
	api.Get("/post/:id", handlers.PostGet)       //  Read
	api.Put("/post/:id", handlers.PostUpdate)    //  Update
	api.Delete("/post/:id", handlers.PostDelete) //  Delete

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
