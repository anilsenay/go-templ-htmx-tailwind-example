package handler

import (
	"errors"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/components"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

type todoRepository interface {
	RetrieveAll() ([]model.Todo, error)
	Insert(todo model.Todo) (model.Todo, error)
	ChangeDone(id int) (*model.Todo, error)
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
	todos, err := h.todoRepository.RetrieveAll()
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}
	return render(ctx, pages.Index(pages.PageProps{Todos: todos}))
}

func (h *TodoHandler) HandlePostTodo(ctx *fiber.Ctx) error {
	todoText := ctx.FormValue("todo-text")
	if todoText == "" {
		return render(ctx, components.AddTodoErrorText(errors.New("Todo text can not be empty")))
	}
	todo, err := h.todoRepository.Insert(model.Todo{Text: todoText})
	if err != nil {
		return render(ctx, components.AddTodoErrorText(errors.New("We cannot complete your request at this time")))
	}
	return combine(
		ctx,
		components.Todo(todo),
		components.AddTodoErrorText(nil),
	)
}

func (h *TodoHandler) HandleUpdateDone(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	todo, err := h.todoRepository.ChangeDone(id)
	if err != nil {
		return err
	}
	return render(ctx, components.Todo(*todo))
}
