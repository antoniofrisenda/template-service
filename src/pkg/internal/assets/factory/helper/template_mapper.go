package helper

import (
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
)

func ToTemplate(d dto.TemplatetPayload) model.Template {
	return model.Template{
		Name:    d.Name,
		Summary: d.Summary,
		Type:    d.Type,
		Content: d.Content,
		Resource: model.Resource{
			URL:       d.Resource.URL,
			Text:      d.Resource.Text,
			Variables: d.Resource.Variables,
		},
	}
}

func ToPayload(m model.Template) *dto.TemplatetPayload {
	return &dto.TemplatetPayload{
		Name:    m.Name,
		Summary: m.Summary,
		Type:    m.Type,
		Content: m.Content,
		Resource: dto.ResourcePayload{
			URL:       m.Resource.URL,
			Text:      m.Resource.Text,
			Variables: m.Resource.Variables,
		},
	}
}
