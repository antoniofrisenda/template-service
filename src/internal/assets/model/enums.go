package model

type DocumentType string

const (
	STATIC   DocumentType = "STATIC"
	TEMPLATE DocumentType = "TEMPLATE"
)

func (e DocumentType) IsValid() bool {
	return e == STATIC || e == TEMPLATE
}

type ContentType string

const (
	PDF        ContentType = "PDF"
	HTML       ContentType = "HTML"
	PLAIN_TEXT ContentType = "PLAIN_TEXT"
	IMAGE      ContentType = "IMAGE"
)

func (e ContentType) IsValid() bool {
	return e == PDF || e == HTML || e == PLAIN_TEXT || e == IMAGE
}

type SourceType string

const (
	FILE SourceType = "FILE"
	TEXT SourceType = "TEXT"
)

func (e SourceType) IsValid() bool {
	return e == FILE || e == TEXT
}
