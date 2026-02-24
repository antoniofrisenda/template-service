package model

type TemplateType string

type ContentType string

const (
	STATIC   TemplateType = "STATIC"
	TEMPLATE TemplateType = "TEMPLATE"
)

const (
	ISPDF   ContentType = "PDF"
	ISHTML  ContentType = "HTML"
	ISPLAIN ContentType = "PLAIN_TEXT"
)

func (t TemplateType) IsValid() bool {
	return t == STATIC || t == TEMPLATE
}

func (c ContentType) IsValid() bool {
	return c == ISPDF || c == ISHTML || c == ISPLAIN
}
