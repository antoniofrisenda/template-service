package router

import (
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service"
	"github.com/gofiber/fiber/v3"
)

type DocumentController struct {
	service service.ITemplateService
}

type IDocumentController interface {
	CreateDocument(c fiber.Ctx) error
	GetDocument(c fiber.Ctx) error
}

func NewDocumentController(service service.ITemplateService) IDocumentController {
	return &DocumentController{
		service: service,
	}
}

// CreateDocument implements [IDocumentController].
func (d *DocumentController) CreateDocument(c fiber.Ctx) error {
	panic("unimplemented")
}

// GetDocument implements [IDocumentController].
func (d *DocumentController) GetDocument(c fiber.Ctx) error {
	panic("unimplemented")
}


