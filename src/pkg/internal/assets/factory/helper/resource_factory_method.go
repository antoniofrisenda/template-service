package helper

import (
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
)

type ResourceMapper interface {
	ToModel(dto.ResourcePayload) (model.Resource, error)
	ToDTO(model.Resource) (dto.ResourcePayload, error)
}

type PlainTextMapper struct{}

func (PlainTextMapper) ToModel(d dto.ResourcePayload) (model.Resource, error) {
	return model.Resource{
		Text:      d.Text,
		Variables: Slice(d.Variables),
	}, nil
}

func (PlainTextMapper) ToDTO(m model.Resource) (dto.ResourcePayload, error) {
	return dto.ResourcePayload{
		Text:      m.Text,
		Variables: Slice(m.Variables),
	}, nil
}

type FileResourceMapper struct{}

func (FileResourceMapper) ToModel(d dto.ResourcePayload) (model.Resource, error) {
	return model.Resource{
		URL:       d.URL,
		Variables: Slice(d.Variables),
	}, nil
}

func (FileResourceMapper) ToDTO(m model.Resource) (dto.ResourcePayload, error) {
	return dto.ResourcePayload{
		URL:       m.URL,
		Variables: Slice(m.Variables),
	}, nil
}

func Slice(src []string) []string {
	if len(src) == 0 {
		return []string{}
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

type Registry struct {
	Map map[model.ContentType]ResourceMapper
}

func NewRegistry() *Registry {
	return &Registry{
		Map: map[model.ContentType]ResourceMapper{
			model.PLAIN_TEXT: PlainTextMapper{},
			model.PDF:        FileResourceMapper{},
			model.HTML:       FileResourceMapper{},
			model.IMAGE:      FileResourceMapper{},
		},
	}
}

func (r *Registry) Get(content model.ContentType) (ResourceMapper, error) {
	return r.Map[content], nil
}
