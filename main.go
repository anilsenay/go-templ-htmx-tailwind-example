package main

import (
	"github.com/a-h/templ"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		return render(c, pages.Index("World!"))
	})

	_ = app.Listen(":8080")
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return component.Render(c.UserContext(), c.Response().BodyWriter())
}
