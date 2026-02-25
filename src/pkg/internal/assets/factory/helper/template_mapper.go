package helper

import (
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITemplateMapper interface {
	ToTemplate(dto.TemplatePayload) (*model.Template, error)
	ToPayload(model.Template) (*dto.TemplatePayload, error)
}

type TemplateMapper struct {
	registry *Registry
}

func NewTemplateMapper(registry *Registry) ITemplateMapper {
	return &TemplateMapper{
		registry: registry,
	}
}

func (tm *TemplateMapper) ToTemplate(d dto.TemplatePayload) (*model.Template, error) {
	mapper, err := tm.registry.Get(d.Content)
	if err != nil {
		return nil, err
	}

	resource, err := mapper.ToModel(d.Resource)
	if err != nil {
		return nil, err
	}

	id := primitive.NilObjectID
	if d.ID != "" {
		id, err = primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return nil, err
		}
	}

	return &model.Template{
		ID:       id,
		Name:     d.Name,
		Summary:  d.Summary,
		Type:     d.Type,
		Content:  d.Content,
		Resource: resource,
	}, nil
}

func (tm *TemplateMapper) ToPayload(m model.Template) (*dto.TemplatePayload, error) {
	mapper, err := tm.registry.Get(m.Content)
	if err != nil {
		return nil, err
	}

	resource, err := mapper.ToDTO(m.Resource)
	if err != nil {
		return nil, err
	}

	return &dto.TemplatePayload{
		ID:       m.ID.Hex(),
		Name:     m.Name,
		Summary:  m.Summary,
		Type:     m.Type,
		Content:  m.Content,
		Resource: resource,
	}, nil
}
