package router

import (
	"encoding/json"

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
	ID := c.Params("ID")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}

	result, err := d.service.FindTemplate(c.Context(), ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Template not found: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (d *documentController) GetPresigned(c fiber.Ctx) error {
	ID := c.Params("ID")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}

	url, err := d.service.FindTemplateWithPresignedURL(c.Context(), ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate presigned URL: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": url,
	})
}

func (d *documentController) PostTemplate(c fiber.Ctx) error {
	payload := dto.InsertDocument{
		Body: &dto.InsertBody{},
	}

	file, err := c.FormFile("file")
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return fiber.NewError(fiber.StatusBadRequest, "File upload error: "+err.Error())
	}

	if file != nil {
		payload.Name = c.FormValue("name")
		payload.Summary = c.FormValue("summary")
		payload.Type = model.DocumentType(c.FormValue("type"))
		payload.Source = model.SourceType("FILE")
		payload.ContentType = model.ContentType(c.FormValue("contentType"))

		result, err := d.service.InsertTemplate(c.Context(), &payload, file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert template: "+err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	}

	if c.Is("json") {
		var temp dto.InsertDocument
		if err := json.Unmarshal(c.Body(), &temp); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		}
		payload = temp

		if err := config.ValidateDocument(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Validation failed: "+err.Error())
		}

		result, err := d.service.InsertTemplate(c, &payload, nil)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to insert template: "+err.Error())
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	}

	return fiber.NewError(fiber.StatusBadRequest, "Request must be either multipart/form-data or application/json")
}

func (d *documentController) GetLatestVariables(c fiber.Ctx) error {
	ID := c.Params("ID")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}

	variables, err := d.service.ExtractVariables(c, ID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Failed to get latest extracted variables: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"variables": variables,
	})
}

func NewDocumentController(service service.DocumentService) DocumentController {
	return &documentController{service: service}
}
