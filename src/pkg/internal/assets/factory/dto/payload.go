package dto

import "github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"

type TemplatetPayload struct {
	Name     string             `json:"name"`
	Summary  string             `json:"summary"`
	Type     model.TemplateType `json:"type"`
	Content  model.ContentType  `json:"content"`
	Resource ResourcePayload    `json:"resource"`
}

type ResourcePayload struct {
	URL       string   `json:"url,omitempty"`
	Text      string   `json:"text,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

type ResolvePayload struct {
	Variables map[string]string `json:"variables"`
}

type S3UploadBase64Payload struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Base64Data  string `json:"base64Data"`
}

type S3UploadBytesPayload struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Bytes       []byte `json:"bytes"`
}
