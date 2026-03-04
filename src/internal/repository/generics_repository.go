package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CRUDRepository[T any] struct {
	collection *mongo.Collection
}

func NewRepository[T any](m *mongo.Collection) *CRUDRepository[T] {
	return &CRUDRepository[T]{collection: m}
}

func (repo *CRUDRepository[T]) Find(ctx context.Context, ID primitive.ObjectID) (*T, error) {
	var t T
	if err := repo.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return &t, nil
}

func (repo *CRUDRepository[T]) Insert(ctx context.Context, t *T) (*T, error) {
	if t == nil {
		return nil, fmt.Errorf("cannot insert nil document")
	}

	_, err := repo.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	return t, nil
}
