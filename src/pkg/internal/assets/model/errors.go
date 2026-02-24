package model

import "errors"

var (
	ErrTemplateNotFound        = errors.New("template not found")
	ErrInvalidTemplateType     = errors.New("invalid template type")
	ErrInvalidContentType      = errors.New("invalid content type")
	ErrInvalidResource         = errors.New("invalid resource")
	ErrStaticCannotHaveVars    = errors.New("static template cannot have variables")
	ErrStaticCannotUsePlain    = errors.New("static template cannot use plain text")
	ErrTemplateMustHaveSource  = errors.New("template must have url or text source")
	ErrTemplateResolveMismatch = errors.New("template variables mismatch")
)
