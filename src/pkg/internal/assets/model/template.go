package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Template struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Summary  string             `bson:"summary"`
	Type     TemplateType       `bson:"type"`
	Content  ContentType        `bson:"content"`
	Resource Resource           `bson:"resource"`
}

type Resource struct {
	URL       string   `bson:"url,omitempty"`
	Text      string   `bson:"text,omitempty"`
	Variables []string `bson:"variables,omitempty"`
}
