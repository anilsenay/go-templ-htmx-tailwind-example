package main

import (
	"context"
	"fmt"

	"github.com/anilsenay/go-htmx-example/handler"
	"github.com/anilsenay/go-htmx-example/repository"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "123456", "go-htmx-example")

	cfg, _ := pgxpool.ParseConfig(psqlInfo)
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	todoRepository := repository.NewTodoRepository(db)
	collectionRepository := repository.NewCollectionRepository(db)
	todoHandler := handler.NewTodoHandler(todoRepository, collectionRepository)
	collectionHandler := handler.NewCollectionHandler(collectionRepository)

	// static files
	app.Static("/", "./public")

	// todo routes
	app.Get("/:collectionId", todoHandler.HandleTodoPage)
	app.Get("/", todoHandler.HandleIndexPage)
	app.Put("/:collectionId/todo/:id/done", todoHandler.HandleUpdateDone)
	app.Post("/:collectionId/todo", todoHandler.HandlePostTodo)
	app.Delete("/:collectionId/todo/:id", todoHandler.HandleDeleteTodo)

	// collection routes
	app.Get("/collection/new", collectionHandler.HandleNewModal)
	app.Get("/collection/edit/:id", collectionHandler.HandleEditModal)
	app.Get("/collection/close", collectionHandler.HandleCloseModal)
	app.Post("/collection/", collectionHandler.HandleCreate)
	app.Put("/collection/:id", collectionHandler.HandleUpdate)
	app.Delete("/collection/:id", collectionHandler.HandleDelete)

	// 404
	app.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		c.Status(fiber.StatusNotFound)
		return pages.NotFound().Render(c.UserContext(), c.Response().BodyWriter())
	})

	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("UI proxy is running on http://localhost:3000")
	_ = app.Listen(":8080")
}
