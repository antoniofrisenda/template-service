package router

import (
	"encoding/json"
	"mime/multipart"
	"strings"

	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
	"github.com/antoniofrisenda/template-service/src/internal/config"
	"github.com/antoniofrisenda/template-service/src/internal/service"
	"github.com/gofiber/fiber/v3"
)

type DocumentController interface {
	GetTemplate(c fiber.Ctx) error
	GetPresigned(c fiber.Ctx) error
	PostTemplate(c fiber.Ctx) error

	GetLatestVariables(c fiber.Ctx) error
}

type documentController struct {
	service service.DocumentService
}

func (d *documentController) GetTemplate(c fiber.Ctx) error {
	id, err := d.getIDParam(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := d.service.FindTemplate(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Template not found: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (d *documentController) GetPresigned(c fiber.Ctx) error {
	id, err := d.getIDParam(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	url, err := d.service.FindTemplateWithPresignedURL(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate presigned URL: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url})
}

func (d *documentController) PostTemplate(c fiber.Ctx) error {
	var (
		payload *dto.InsertDocument
		file    *multipart.FileHeader
		err     error
	)

	header := c.Get("Content-Type")

	switch {
	case header != "" && strings.HasPrefix(header, "multipart/form-data"):
		payload, file, err = d.parseMultipart(c)
	case header == "application/json":
		payload, err = d.parseJSON(c)
	default:
		return fiber.NewError(fiber.StatusBadRequest, "Unsupported content type")
	}

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := d.service.InsertTemplate(c.Context(), payload, file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (d *documentController) GetLatestVariables(c fiber.Ctx) error {
	id, err := d.getIDParam(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	variables, err := d.service.ExtractVariables(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"variables": variables})
}

func NewDocumentController(service service.DocumentService) DocumentController {
	return &documentController{service: service}
}

func (d *documentController) getIDParam(c fiber.Ctx) (string, error) {
	ID := c.Params("ID")
	if ID == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}
	return ID, nil
}

func (d *documentController) parseMultipart(c fiber.Ctx) (*dto.InsertDocument, *multipart.FileHeader, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "File upload error: "+err.Error())
	}

	return &dto.InsertDocument{
		Name:        c.FormValue("name"),
		Summary:     c.FormValue("summary"),
		Type:        model.DocumentType(c.Params("DocumentType")),
		Source:      model.SourceType("FILE"),
		ContentType: model.ContentType(c.FormValue("contentType")),
		Body:        &dto.InsertBody{},
	}, file, nil
}

func (d *documentController) parseJSON(c fiber.Ctx) (*dto.InsertDocument, error) {
	var payload dto.InsertDocument
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid JSON payload: "+err.Error())
	}

	if err := config.Validate(&payload); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return &payload, nil
}
