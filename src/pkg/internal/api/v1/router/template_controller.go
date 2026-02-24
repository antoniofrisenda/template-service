package router

import (
	"strings"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/gofiber/fiber/v3"
)

type ITemplateController interface {
	Find(c fiber.Ctx) error
	Create(c fiber.Ctx) error

	Delete(c fiber.Ctx) error
	Patch(c fiber.Ctx) error

	Upload(c fiber.Ctx, data interface{}, contentType string) error
	Download(c fiber.Ctx) error

	PresignedURL(c fiber.Ctx) error
}

type TemplateController struct {
	service service.ITemplateService
}

func NewTemplateController(service service.ITemplateService) ITemplateController {
	return &TemplateController{
		service: service,
	}
}

func (tc *TemplateController) Find(c fiber.Ctx) error {
	ID := c.Params("id")

	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing param=%s:", ID)
	}

	t, err := tc.service.Find(c, ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(t)
}

func (tc *TemplateController) Create(c fiber.Ctx) error {
	var payload dto.TemplatetPayload
	if err := c.Bind().Body(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	created, err := tc.service.Create(c, payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (tc *TemplateController) Delete(c fiber.Ctx) error {
	ID := c.Params("id")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing param=%s:", ID)
	}

	ok, err := tc.service.Delete(c, ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted": ok})
}

func (tc *TemplateController) Patch(c fiber.Ctx) error {
	ID := c.Params("id")
	if ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing param=%s:", ID)
	}

	var payload dto.TemplatetPayload
	if err := c.Bind().Body(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	updated, err := tc.service.Patch(c, ID, payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (tc *TemplateController) Upload(c fiber.Ctx, data interface{}, contentType string) error {
	key := c.Params("key")

	switch d := data.(type) {
	case string:
		if err := tc.service.UploadBase64(c, key, d, contentType); err != nil {
			return err
		}
	case []byte:
		if err := tc.service.UploadBytes(c, key, d, contentType); err != nil {
			return err
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported upload data type")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"uploaded": true,
		"key":      key,
	})
}

func (tc *TemplateController) Download(c fiber.Ctx) error {
	key := c.Params("key")

	format := strings.ToLower(c.Query("format"))

	accept := strings.ToLower(c.Get("Accept"))
	if format == "" && strings.Contains(accept, "application/octet-stream") {
		format = "bytes"
	}
	if format == "" {
		format = "base64"
	}

	switch format {
	case "bytes":
		data, err := tc.service.DownloadBytes(c, key)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "application/octet-stream")
		return c.Status(fiber.StatusOK).Send(data)

	case "base64":
		data, err := tc.service.DownloadBase64(c, key)
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/plain")
		return c.Status(fiber.StatusOK).SendString(data)
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid format")
	}
}

func (tc *TemplateController) PresignedURL(c fiber.Ctx) error {
	url, err := tc.service.GetPresignedURL(c, c.Params("key"))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url})
}
