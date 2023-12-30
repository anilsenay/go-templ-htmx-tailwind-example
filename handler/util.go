package handler

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return component.Render(c.UserContext(), c.Response().BodyWriter())
}
