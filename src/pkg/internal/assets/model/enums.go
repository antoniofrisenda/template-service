package model

type DocumentType string
type ContentType string
type Source string

const (
	STATIC   DocumentType = "STATIC"
	TEMPLATE DocumentType = "TEMPLATE"
)

const (
	S3_PRIVATE_BUCKET Source = "S3_PRIVATE_BUCKET"
	HARDCODED_VALUE Source = "HARDCODED_VALUE"
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
