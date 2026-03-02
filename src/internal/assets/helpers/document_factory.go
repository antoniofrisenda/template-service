package helpers

import (
	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
)

func Register(dto *dto.InsertDocument) *model.Document {
	if dto == nil {
		return nil
	}

	contentType := model.ContentType(dto.ContentType)

	switch dto.Type {
	case model.STATIC:
		switch dto.Source {
		case "TEXT":
			if dto.Body == nil || dto.Body.Text == nil {
				return nil
			}
			return model.NewStaticTextDocument(dto.Name, dto.Summary, contentType, *dto.Body.Text)
		case "FILE":
			return model.NewStaticFileDocument(dto.Name, dto.Summary, contentType, "")
		}
	case model.TEMPLATE:
		switch dto.Source {
		case "TEXT":
			if dto.Body == nil || dto.Body.Text == nil {
				return nil
			}
			return model.NewTemplateTextDocument(dto.Name, dto.Summary, contentType, *dto.Body.Text, dto.Body.Variables)
		case "FILE":
			return model.NewTemplateFileDocument(dto.Name, dto.Summary, contentType, "", dto.Body.Variables)
		}
	}

	return nil
}
