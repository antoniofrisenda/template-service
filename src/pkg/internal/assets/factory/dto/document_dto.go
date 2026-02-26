type DocumentDto struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Name    string `bson:"name"`
	Summary string `bson:"summary,omitempty"`

	Content ContentType  `bson:"content"`

	Resource *ResourceDto `bson:"resource,omitempty"`
	Template *TemplateDto `bson:"template,omitempty"`
}

type ResourceDto struct {
	URL  *string `bson:"url,omitempty"`
	Text *string `bson:"text,omitempty"`
}

type TemplateDto struct {
	URL       *string  `bson:"url,omitempty"`
	Text      *string  `bson:"text,omitempty"`
	Variables []string `bson:"variables,omitempty"`
}