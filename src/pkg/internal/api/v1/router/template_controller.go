package router

import (
	"fmt"
	"io"
	"strings"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
	"github.com/gofiber/fiber/v3"
)

type ITemplateController interface {
	FindTemplate(c fiber.Ctx) error

	FindTemplateByName(c fiber.Ctx) error

	FindTemplateBySummary(c fiber.Ctx) error
	CreateTemplate(c fiber.Ctx) error

	DeleteTemplate(c fiber.Ctx) error
	PatchTemplate(c fiber.Ctx) error

	UploadFile(c fiber.Ctx) error

	DownloadFileByBase64(c fiber.Ctx) error

	DownloadFileByBytes(c fiber.Ctx) error

	DownloadFileByPresignedURL(c fiber.Ctx) error
}

type TemplateController struct {
	service service.ITemplateService
}

func NewTemplateController(service service.ITemplateService) ITemplateController {
	return &TemplateController{
		service: service,
	}
}

func (tc *TemplateController) FindTemplate(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing param=%s:", id)
	}

	find, err := tc.service.Find(c, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to find template=%v: %s", find, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(find)
}

func (tc *TemplateController) FindTemplateByName(c fiber.Ctx) error {
	var name string
	find, err := tc.service.SearchTemplateName(c, name)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to find template=%v: %s", find, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(find)
}

func (tc *TemplateController) FindTemplateBySummary(c fiber.Ctx) error {
	var period string
	find, err := tc.service.Find(c, period)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to find template=%v: %s", find, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(find)
}

func (tc *TemplateController) CreateTemplate(c fiber.Ctx) error {
	var payload dto.TemplatePayload
	if err := c.Bind().Body(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	create, err := tc.service.Create(c, payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("payload provided=%v: %s", payload, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(create)
}

func (tc *TemplateController) PatchTemplate(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing param=%s:", id)
	}

	var payload dto.TemplatePayload
	if err := c.Bind().Body(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	update, err := tc.service.Patch(c, id, payload)
	if err != nil {
		return fiber.NewError(fiber.StatusNotModified, fmt.Sprintf("failed to update template=%v: %s", id, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(update)
}

func (tc *TemplateController) DeleteTemplate(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("missing param=%s: ", id))
	}

	delete, err := tc.service.Delete(c, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotModified, fmt.Sprintf("failed to delete template=%v: %s", id, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"deleted": delete})
}

func (tc *TemplateController) UploadFile(c fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "missing file form field")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to open uploaded file: %s", err.Error()))
	}
	defer file.Close()

	input, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to read uploaded file: %s", err.Error()))
	}

	content := model.ContentType(strings.ToUpper(strings.TrimSpace(c.FormValue("content"))))
	templateType := model.TemplateType(strings.ToUpper(strings.TrimSpace(c.FormValue("type"))))
	name := strings.TrimSpace(c.FormValue("name"))
	summary := strings.TrimSpace(c.FormValue("summary"))

	if name == "" || summary == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name and summary are required")
	}

	if !content.IsValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid content, accepted values: HTML, PDF, PLAIN_TEXT")
	}

	if !templateType.IsValid() {
		return fiber.NewError(fiber.StatusBadRequest, "invalid type, accepted values: TEMPLATE, STATIC")
	}

	create, err := tc.service.CreateByUploadingFile(c, input, name, summary, templateType, content)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to upload file: %s", err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(create)
}

func (tc *TemplateController) DownloadFileByBase64(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("missing param=%s: ", id))
	}

	data, err := tc.service.DownloadBase64(c, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to download template=%v: %s", id, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

func (tc *TemplateController) DownloadFileByBytes(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("missing param=%s: ", id))
	}

	data, err := tc.service.DownloadBytes(c, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to download template=%v: %s", id, err.Error()))
	}

	return c.Status(fiber.StatusOK).Send(data)
}

func (tc *TemplateController) DownloadFileByPresignedURL(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("missing param=%s: ", id))
	}

	url, err := tc.service.DownloadByPresignedURL(c, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to generate presigned url template=%v: %s", id, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url})
}
