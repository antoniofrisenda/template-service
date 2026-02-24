package middleware

import (
	"errors"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := err.Error()

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code
		message = fiberErr.Message
	} else {
		switch {
		case errors.Is(err, model.ErrTemplateNotFound), errors.Is(err, mongo.ErrNoDocuments):
			status = fiber.StatusNotFound
		case errors.Is(err, model.ErrInvalidTemplateType),
			errors.Is(err, model.ErrInvalidContentType),
			errors.Is(err, model.ErrInvalidResource),
			errors.Is(err, model.ErrStaticCannotHaveVars),
			errors.Is(err, model.ErrStaticCannotUsePlain),
			errors.Is(err, model.ErrTemplateMustHaveSource),
			errors.Is(err, model.ErrTemplateResolveMismatch):
			status = fiber.StatusBadRequest
		}
	}

	return c.Status(status).JSON(fiber.Map{"error": message})
}
