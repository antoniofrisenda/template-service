package helper

import (
	"net/url"
	"strings"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
)

func ValidateTemplate(tpl model.Template) error {
	name := strings.TrimSpace(tpl.Name)
	summary := strings.TrimSpace(tpl.Summary)
	text := strings.TrimSpace(tpl.Resource.Text)
	urlStr := strings.TrimSpace(tpl.Resource.URL)

	if !tpl.Type.IsValid() {
		return model.ErrInvalidTemplateType
	}
	if !tpl.Content.IsValid() {
		return model.ErrInvalidContentType
	}

	if name == "" {
		return model.ErrInvalidField
	}
	if summary == "" {
		return model.ErrInvalidField
	}

	switch {
	case tpl.Type == model.STATIC && tpl.Content == model.PLAIN_TEXT:
		if len(tpl.Resource.Variables) > 0 {
			return model.ErrStaticCannotHaveVars
		}

		if urlStr != "" {
			return model.ErrStaticCannotHaveURL
		}

	case tpl.Type == model.STATIC && tpl.Content == model.IMAGE:
		if len(tpl.Resource.Variables) > 0 {
			return model.ErrStaticCannotHaveVars
		}
		if urlStr == "" {
			return model.ErrInvalidField
		}
		if text != "" {
			return model.ErrInvalidField
		}

		u, err := url.ParseRequestURI(urlStr)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return model.ErrInvalidField
		}

	case tpl.Type == model.TEMPLATE && tpl.Content == model.PLAIN_TEXT:
		if text == "" {
			return model.ErrTemplateMustHaveSource
		}
		if urlStr != "" {
			return model.ErrTemplateCannotHaveURL
		}

	case tpl.Type == model.TEMPLATE &&
		(tpl.Content == model.PDF || tpl.Content == model.HTML):

		if urlStr == "" {
			return model.ErrTemplateMustHaveSource
		}
		if text != "" {
			return model.ErrTemplateCannotHaveText
		}

		u, err := url.ParseRequestURI(urlStr)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return model.ErrInvalidField
		}

	default:
		return model.ErrInvalidTemplateType
	}

	return nil
}
