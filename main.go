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
	app.Get("/", todoHandler.HandleTodoPage)
	app.Put("/todos/:id/done", todoHandler.HandleUpdateDone)
	app.Post("/todos", todoHandler.HandlePostTodo)

	// 404
	app.Get("*", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		c.Status(fiber.StatusNotFound)
		return pages.NotFound().Render(c.UserContext(), c.Response().BodyWriter())
	})

	_ = app.Listen(":8080")
}
