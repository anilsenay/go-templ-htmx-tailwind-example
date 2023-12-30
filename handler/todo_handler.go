package handler

import (
	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

type todoRepository interface {
	RetrieveAll() ([]model.Todo, error)
	Insert(todo model.Todo) error
	SetDone(id int) error
}

type TodoHandler struct {
	todoRepository todoRepository
}

func NewTodoHandler(r todoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepository: r,
	}
}

func (h *TodoHandler) HandleTodoPage(ctx *fiber.Ctx) error {
	todos, _ := h.todoRepository.RetrieveAll()
	return render(ctx, pages.Index(pages.PageProps{Todos: todos}))
}
