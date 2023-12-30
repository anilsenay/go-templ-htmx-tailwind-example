package main

import (
	"github.com/a-h/templ"
	"github.com/anilsenay/go-htmx-example/models"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

var todos = []models.Todo{
	{Text: "Todo1", Done: true},
	{Text: "Todo2"},
}

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, pages.Index(pages.PageProps{Todos: todos}))
	})

	_ = app.Listen(":8080")
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return component.Render(c.UserContext(), c.Response().BodyWriter())
}
