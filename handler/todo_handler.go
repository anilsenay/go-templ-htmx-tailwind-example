package handler

import (
	"errors"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/components"
	"github.com/anilsenay/go-htmx-example/view/layout"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type todoRepository interface {
	GetTodoByCollection(c model.Collection) (model.CollectionWithTodoList, error)
	Insert(collectionId int, todo model.Todo) (model.Todo, error)
	ChangeDone(collectionId int, id int) (*model.Todo, error)
	Delete(collectionId int, id int) error
}

type todoHandlerCollectionRepository interface {
	GetCollections() ([]model.Collection, error)
	GetCollection(id int) (model.Collection, error)
}

type TodoHandler struct {
	todoRepository       todoRepository
	collectionRepository todoHandlerCollectionRepository
}

func NewTodoHandler(tr todoRepository, cr todoHandlerCollectionRepository) *TodoHandler {
	return &TodoHandler{
		todoRepository:       tr,
		collectionRepository: cr,
	}
}

func (h *TodoHandler) HandleIndexPage(ctx *fiber.Ctx) error {
	collections, err := h.collectionRepository.GetCollections()
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

	collections, err := h.collectionRepository.GetCollections()
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	collection, err := h.collectionRepository.GetCollection(collectionId)
	if err != nil {
		return err
	}

	todo, err := h.todoRepository.GetTodoByCollection(collection)
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

	todoText := utils.CopyString(ctx.FormValue("todo-text"))
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
