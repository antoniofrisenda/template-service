package dto

import (
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
)

type Document struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Summary       string             `json:"summary"`
	Type          model.DocumentType `json:"type"`
	Source        model.SourceType   `json:"source"`
	ContentType   model.ContentType  `json:"contentType"`
	Base64Encoded bool               `json:"base64Encoded"`
	Body          string             `json:"body"`
}

type InsertDocument struct {
	Name        string             `json:"name"`
	Summary     string             `json:"summary"`
	Type        model.DocumentType `json:"type"`
	Source      model.SourceType   `json:"source"`
	ContentType model.ContentType  `json:"contentType"`
	Body        *InsertBody        `json:"body"`
}

type InsertBody struct {
	URL       *string  `json:"url,omitempty"`
	Text      *string  `json:"text,omitempty"`
	Variables []string `json:"variables,omitempty"`
}
