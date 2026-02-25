package model

import "errors"

var (
	ErrTemplateNotFound = errors.New("TEMPLATE or STATIC file not found")
	ErrInvalidTemplateType = errors.New("invalid template type")
	ErrInvalidContentType  = errors.New("invalid content type")

	ErrInvalidField = errors.New("template field has invalid syntax or is empty")

	ErrStaticMustHaveText   = errors.New("static template must contain text")
	ErrStaticCannotHaveURL  = errors.New("static template cannot contain a URL")
	ErrStaticCannotHaveVars = errors.New("static template cannot contain variables")

	ErrTemplateMustHaveSource = errors.New("template must contain a valid source (text or URL)")
	ErrTemplateCannotHaveURL  = errors.New("template with plain text content cannot contain a URL")
	ErrTemplateCannotHaveText = errors.New("template with non-text content cannot contain inline text")
)
