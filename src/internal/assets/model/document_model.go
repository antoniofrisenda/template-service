package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Summary     string             `bson:"summary"`
	Type        DocumentType       `bson:"type"`
	Source      SourceType         `bson:"source"`
	ContentType ContentType        `bson:"contentType"`
	Body        *DocumentBody      `bson:"body"`
}

type DocumentBody struct {
	URL       *string  `bson:"url,omitempty"`
	Text      *string  `bson:"text,omitempty"`
	Variables []string `bson:"variables,omitempty"`
}

func NewStaticFileDocument(name string, summary string, contentType ContentType, url string) *Document {
	return &Document{
		Name:        name,
		Summary:     summary,
		Type:        STATIC,
		Source:      FILE,
		ContentType: contentType,
		Body: &DocumentBody{
			URL: &url,
		},
	}
}

func NewStaticTextDocument(name, summary string, contentType ContentType, text string) *Document {
	return &Document{
		Name:        name,
		Summary:     summary,
		Type:        STATIC,
		Source:      TEXT,
		ContentType: contentType,
		Body: &DocumentBody{
			Text: &text,
		},
	}
}

func NewTemplateFileDocument(name, summary string, contentType ContentType, url string, variables []string) *Document {
	return &Document{
		Name:        name,
		Summary:     summary,
		Type:        TEMPLATE,
		Source:      FILE,
		ContentType: contentType,
		Body: &DocumentBody{
			URL:       &url,
			Variables: variables,
		},
	}
}

func NewTemplateTextDocument(name, summary string, contentType ContentType, text string, variables []string) *Document {
	return &Document{
		Name:        name,
		Summary:     summary,
		Type:        TEMPLATE,
		Source:      TEXT,
		ContentType: contentType,
		Body: &DocumentBody{
			Text:      &text,
			Variables: variables,
		},
	}
}
