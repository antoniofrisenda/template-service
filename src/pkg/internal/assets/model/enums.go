package model

type TemplateType string

type ContentType string

const (
	STATIC   TemplateType = "STATIC"
	TEMPLATE TemplateType = "TEMPLATE"
)

const (
	PDF        ContentType = "PDF"
	HTML       ContentType = "HTML"
	PLAIN_TEXT ContentType = "PLAIN_TEXT"

	IMAGE ContentType = "IMAGE"
)

func (t TemplateType) IsValid() bool {
	return t == STATIC || t == TEMPLATE
}

func (c ContentType) IsValid() bool {
	return c == PDF || c == HTML || c == PLAIN_TEXT || c == IMAGE
}
