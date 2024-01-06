package handler

import (
	"fmt"

	"github.com/anilsenay/go-htmx-example/model"
	"github.com/anilsenay/go-htmx-example/view/components"
	"github.com/anilsenay/go-htmx-example/view/pages"
	"github.com/gofiber/fiber/v2"
)

type collectionRepository interface {
	GetCollection(id int) (model.Collection, error)
	InsertCollection(collection model.Collection) (model.Collection, error)
	UpdateCollection(id int, updates model.Collection) (model.Collection, error)
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

func (h *CollectionHandler) HandleCreateCollection(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name")
	color := ctx.FormValue("color")
	collection, err := h.collectionRepository.InsertCollection(model.Collection{Name: name, HexColor: color})
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	ctx.Set("HX-Location", fmt.Sprintf("/%d", collection.Id))
	return nil
}

func (h *CollectionHandler) HandleUpdateCollection(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return err
	}

	name := ctx.FormValue("name")
	color := ctx.FormValue("color")
	collection, err := h.collectionRepository.UpdateCollection(id, model.Collection{Name: name, HexColor: color})
	if err != nil {
		return render(ctx, pages.ErrorPage(fiber.StatusInternalServerError, "Internal Server Error: "+err.Error()))
	}

	ctx.Set("HX-Location", fmt.Sprintf("/%d", collection.Id))
	return nil
}
