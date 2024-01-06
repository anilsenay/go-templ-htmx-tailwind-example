package handler

import (
	"errors"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/components"
	"github.com/anilsenay/go-htmx-example/view/layout"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

type todoRepository interface {
	GetCollections() ([]model.Collection, error)
	GetTodoByCollection(id int) (model.CollectionWithTodoList, error)
	Insert(collectionId int, todo model.Todo) (model.Todo, error)
	ChangeDone(collectionId int, id int) (*model.Todo, error)
	Delete(collectionId int, id int) error
}

type TodoHandler struct {
	todoRepository todoRepository
}

func NewTodoHandler(r todoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepository: r,
	}
}

func (h *TodoHandler) HandleIndexPage(ctx *fiber.Ctx) error {
	collections, err := h.todoRepository.GetCollections()
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}
	return render(ctx, pages.TodoPage(pages.TodoPageProps{
		PageLayoutProps: layout.PageProps{MenuItems: collections, Title: "My Collections"},
	}))
}

func (h *TodoHandler) HandleTodoPage(ctx *fiber.Ctx) error {
	collectionId, err := ctx.ParamsInt("collectionId")
	if err != nil {
		return err
	}

	collections, err := h.todoRepository.GetCollections()
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	todo, err := h.todoRepository.GetTodoByCollection(collectionId)
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}
	return render(ctx, pages.TodoPage(pages.TodoPageProps{
		CollectionWithTodo: todo,
		PageLayoutProps:    layout.PageProps{Active: collectionId, MenuItems: collections, Title: todo.Name},
	}))
}

func (h *TodoHandler) HandlePostTodo(ctx *fiber.Ctx) error {
	collectionId, err := ctx.ParamsInt("collectionId")
	if err != nil {
		return err
	}

	todoText := ctx.FormValue("todo-text")
	if todoText == "" {
		return render(ctx, components.AddTodoErrorText(errors.New("Todo text can not be empty")))
	}
	todo, err := h.todoRepository.Insert(collectionId, model.Todo{Text: todoText})
	if err != nil {
		return render(ctx, components.AddTodoErrorText(errors.New("We cannot complete your request at this time")))
	}
	return combine(
		ctx,
		components.Todo(collectionId, todo),
		components.AddTodoErrorText(nil),
	)
}

func (h *TodoHandler) HandleUpdateDone(ctx *fiber.Ctx) error {
	collectionId, err := ctx.ParamsInt("collectionId")
	if err != nil {
		return err
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	todo, err := h.todoRepository.ChangeDone(collectionId, id)
	if err != nil {
		return err
	}
	return render(ctx, components.Todo(collectionId, *todo))
}

func (h *TodoHandler) HandleDeleteTodo(ctx *fiber.Ctx) error {
	collectionId, err := ctx.ParamsInt("collectionId")
	if err != nil {
		return err
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	err = h.todoRepository.Delete(collectionId, id)
	if err != nil {
		return err
	}
	return nil
}
