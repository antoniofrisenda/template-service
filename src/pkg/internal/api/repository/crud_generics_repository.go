package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CRUDRepository[T any] struct {
	collection *mongo.Collection
}

func NewGenericRepository[T any](m *mongo.Collection) *CRUDRepository[T] {
	return &CRUDRepository[T]{collection: m}
}

func (repo *CRUDRepository[T]) Create(ctx context.Context, t *T) (primitive.ObjectID, error) {
	query, err := repo.collection.InsertOne(ctx, t)
	if err != nil {
		return primitive.NilObjectID, err
	}
	ID, ok := query.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, err
	}
	return ID, nil
}

func (repo *CRUDRepository[T]) FindByID(ctx context.Context, ID primitive.ObjectID) (*T, error) {
	var t T
	if err := repo.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *CRUDRepository[T]) UpdateByID(ctx context.Context, ID primitive.ObjectID, b bson.M) error {
	_, err := repo.collection.UpdateByID(ctx, ID, bson.M{"$set": b})
	return err
}

func (repo *CRUDRepository[T]) DeleteByID(ctx context.Context, ID primitive.ObjectID) error {
	_, err := repo.collection.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}
