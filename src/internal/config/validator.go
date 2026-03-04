package config

import (
	"fmt"
	"strings"

	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
)

type Validator interface {
	Validate() error
}

func Validate(d *dto.InsertDocument) error {
	if d == nil {
		return fmt.Errorf("document is nil")
	}

	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if !d.Type.IsValid() {
		return fmt.Errorf("invalid document type: %s (must be STATIC or TEMPLATE)", d.Type)
	}

	if !d.Source.IsValid() {
		return fmt.Errorf("invalid source type: %s (must be FILE or TEXT)", d.Source)
	}

	if !d.ContentType.IsValid() {
		return fmt.Errorf("invalid content type: %s (must be PDF, HTML, PLAIN_TEXT, or IMAGE)", d.ContentType)
	}

	if d.ContentType == model.IMAGE {
		if d.Source != model.FILE {
			return fmt.Errorf("IMAGE content type requires FILE source")
		}
		if d.Type == model.TEMPLATE {
			return fmt.Errorf("IMAGE content type cannot be used with TEMPLATE documents")
		}
	}

	if d.Source == model.TEXT {
		if d.Body == nil || d.Body.Text == nil || strings.TrimSpace(*d.Body.Text) == "" {
			return fmt.Errorf("text source requires non-empty text in body")
		}
	}

	return nil
}
