package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm-crud/database"
	"gorm-crud/handlers"
	"gorm-crud/middlewares"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{AppName: "Gorm Crud Post App"})
	app.Use(logger.New())
	database.ConnectDB()
	api := app.Group("/api")

	api.Post("/signup", handlers.SignUp)
	api.Post("/signin", handlers.SignIn)

	auth := middlewares.AuthMiddleware()
	api.Use(auth)

	api.Post("/post", handlers.PostCreate)       //  Create
	api.Get("/post/:id", handlers.PostGet)       //  Read
	api.Put("/post/:id", handlers.PostUpdate)    //  Update
	api.Delete("/post/:id", handlers.PostDelete) //  Delete

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
