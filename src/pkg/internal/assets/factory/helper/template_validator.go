package helper

import "github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"

func ValidateTemplate(tpl model.Template) error {
	if !tpl.Type.IsValid() {
		return model.ErrInvalidTemplateType
	}
	if !tpl.Content.IsValid() {
		return model.ErrInvalidContentType
	}

	if tpl.Type == model.STATIC {
		if tpl.Content == model.ISPLAIN {
			return model.ErrStaticCannotUsePlain
		}
		if len(tpl.Resource.Variables) > 0 {
			return model.ErrStaticCannotHaveVars
		}
		if tpl.Resource.URL == "" {
			return model.ErrInvalidResource
		}
		return nil
	}

	if tpl.Type == model.TEMPLATE {
		if tpl.Content == model.ISPLAIN && tpl.Resource.Text == "" {
			return model.ErrTemplateMustHaveSource
		}
		if tpl.Content != model.ISPLAIN && tpl.Resource.URL == "" && tpl.Resource.Text == "" {
			return model.ErrTemplateMustHaveSource
		}
		return nil
	}

	return model.ErrInvalidTemplateType
}
