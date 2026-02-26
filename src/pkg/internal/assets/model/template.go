package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name    string `bson:"name"`
	Summary string `bson:"summary,omitempty"`

	Type    TemplateType `bson:"type"`
	Source  Source       `bson:"source"`
	Content ContentType  `bson:"content"`

	Resource *Resource `bson:"resource,omitempty"`
	Template *Template `bson:"template,omitempty"`
}

type Resource struct {
	URL  *string `bson:"url,omitempty"`
	Text *string `bson:"text,omitempty"`
}

type Template struct {
	URL       *string  `bson:"url,omitempty"`
	Text      *string  `bson:"text,omitempty"`
	Variables []string `bson:"variables,omitempty"`
}

func NewInternalStaticDocument(
	code string,
	name string,
	summary string,
	content ContentType,
	url string,
) *Document {

	return &Document{
		Name:    name,
		Summary: summary,
		Type:    STATIC,
		Source:  S3_PRIVATE_BUCKET,
		Content: content,
		Resource: &Resource{
			URL: &url,
		},
	}
}

func NewInternalTemplateDocument(
	code string,
	name string,
	summary string,
	content ContentType,
	url string,
	variables []string,
) *Document {

	return &Document{
		Name:    name,
		Summary: summary,
		Type:    TEMPLATE,
		Source:  S3_PRIVATE_BUCKET,
		Content: content,
		Template: &Template{
			URL:       &url,
			Variables: variables,
		},
	}
}


func NewHardcodedStaticDocument(
	code string,
	name string,
	summary string,
	content ContentType,
	text string,
) *Document {

	return &Document{
		Name:    name,
		Summary: summary,
		Type:    STATIC,
		Source:  HARDCODED_VALUE,
		Content: content,
		Resource: &Resource{
			Text: &text,
		},
	}
}


func NewHardcodedTemplateDocument(
	code string,
	name string,
	summary string,
	content ContentType,
	text string,
	variables []string,
) *Document {

	return &Document{
		Name:    name,
		Summary: summary,
		Type:    TEMPLATE,
		Source:  HARDCODED_VALUE,
		Content: content,
		Template: &Template{
			Text:      &text,
			Variables: variables,
		},
	}
}
