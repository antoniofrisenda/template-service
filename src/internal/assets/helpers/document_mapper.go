package helpers

import (
	"fmt"

	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
)

type DocumentMapper interface {
	ToDTO(m *model.Document) (*dto.Document, error)
	ToModel(m *dto.InsertDocument) (*model.Document, error)
}

type documentMapper struct{}

func NewDocumentMapper() DocumentMapper {
	return &documentMapper{}
}

func (dm *documentMapper) ToDTO(m *model.Document) (*dto.Document, error) {
	var (
		base64Encoded bool
		body          string
	)

	if m == nil {
		return nil, fmt.Errorf("document is nil")
	}

	if !m.Type.IsValid() {
		return nil, fmt.Errorf("invalid document type: %s", m.Type)
	}

	if !m.ContentType.IsValid() {
		return nil, fmt.Errorf("invalid content type: %s", m.ContentType)
	}

	if !m.Source.IsValid() {
		return nil, fmt.Errorf("invalid source type: %s", m.Source)
	}

	switch m.Source {
	case model.FILE:
		base64Encoded = true

		if m.Body == nil || m.Body.URL == nil {
			return nil, fmt.Errorf("file source requires URL in body")
		}

		body = *m.Body.URL

	case model.TEXT:
		base64Encoded = false

		if m.Body == nil || m.Body.Text == nil {
			return nil, fmt.Errorf("text source requires text in body")
		}

		body = *m.Body.Text
	default:
		return nil, fmt.Errorf("unsupported source type: %s", m.Source)
	}

	return &dto.Document{
		ID:            m.ID.Hex(),
		Name:          m.Name,
		Summary:       m.Summary,
		Type:          m.Type,
		Source:        m.Source,
		ContentType:   m.ContentType,
		Base64Encoded: base64Encoded,
		Body:          body,
	}, nil
}

func (dm *documentMapper) ToModel(d *dto.InsertDocument) (*model.Document, error) {
	return Register(d), nil
}
