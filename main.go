package main

import (
	"github.com/anilsenay/go-htmx-example/handler"
	"github.com/anilsenay/go-htmx-example/repository"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})

	// static files
	app.Static("/", "./public")

	todoRepository := repository.NewTodoRepository()
	todoHandler := handler.NewTodoHandler(todoRepository)

	// todo routes
	app.Get("/:collectionId", todoHandler.HandleTodoPage)
	app.Get("/", todoHandler.HandleIndexPage)
	app.Put("/:collectionId/todo/:id/done", todoHandler.HandleUpdateDone)
	app.Post("/:collectionId/todo", todoHandler.HandlePostTodo)
	app.Delete("/:collectionId/todo/:id", todoHandler.HandleDeleteTodo)

	// 404
	app.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		c.Status(fiber.StatusNotFound)
		return pages.NotFound().Render(c.UserContext(), c.Response().BodyWriter())
	})

	_ = app.Listen(":8080")
}
