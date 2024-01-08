package handler

import (
	"fmt"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/components"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type collectionRepository interface {
	GetCollection(id int) (model.Collection, error)
	Insert(collection model.Collection) (model.Collection, error)
	Update(id int, updates model.Collection) (model.Collection, error)
	Delete(id int) error
}

type CollectionHandler struct {
	collectionRepository collectionRepository
}

func NewCollectionHandler(r collectionRepository) *CollectionHandler {
	return &CollectionHandler{
		collectionRepository: r,
	}
}

func (h *CollectionHandler) HandleNewModal(ctx *fiber.Ctx) error {
	return render(ctx, components.AddCollectionModal())
}

func (h *CollectionHandler) HandleEditModal(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}
	collection, err := h.collectionRepository.GetCollection(id)
	if err != nil {
		return err
	}

	return render(ctx, components.EditCollectionModal(collection))
}

func (h *CollectionHandler) HandleCloseModal(ctx *fiber.Ctx) error {
	return nil
}

func (h *CollectionHandler) HandleCreate(ctx *fiber.Ctx) error {
	name := utils.CopyString(ctx.FormValue("name"))
	color := utils.CopyString(ctx.FormValue("color"))
	collection, err := h.collectionRepository.Insert(model.Collection{Name: name, Color: color})
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	ctx.Set("HX-Location", fmt.Sprintf("/%d", collection.Id))
	return nil
}

func (h *CollectionHandler) HandleUpdate(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	name := utils.CopyString(ctx.FormValue("name"))
	color := utils.CopyString(ctx.FormValue("color"))
	collection, err := h.collectionRepository.Update(id, model.Collection{Name: name, Color: color})
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	return combine(
		ctx,
		components.Collection(collection, true),
		components.TodoTitle(collection),
	)
}

func (h *CollectionHandler) HandleDelete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	err = h.collectionRepository.Delete(id)
	if err != nil {
		return err
	}

	ctx.Set("HX-Location", "/")
	return nil
}
