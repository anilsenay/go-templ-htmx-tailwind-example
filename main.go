package main

import (
	"github.com/anilsenay/go-htmx-example/handler"
	"github.com/anilsenay/go-htmx-example/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Static("/", "./public")

	todoRepository := repository.NewTodoRepository()
	todoHandler := handler.NewTodoHandler(todoRepository)
	app.Get("/", todoHandler.HandleTodoPage)
	app.Put("/todos/:id/done", todoHandler.HandleUpdateDone)

	_ = app.Listen(":8080")
}
