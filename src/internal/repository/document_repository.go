package repository

import (
	"context"

	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentRepository interface {
	FindOne(ctx context.Context, id primitive.ObjectID) (*model.Document, error)
	InsertOne(ctx context.Context, m *model.Document) (*model.Document, error)
	EnsureIndexes(ctx context.Context) error
}

type documentRepository struct {
	repo       *CRUDRepository[model.Document]
	collection *mongo.Collection
}

func NewDocumentRepository(collection *mongo.Collection) DocumentRepository {
	return &documentRepository{
		repo:       NewRepository[model.Document](collection),
		collection: collection,
	}
}

func (r *documentRepository) FindOne(ctx context.Context, ID primitive.ObjectID) (*model.Document, error) {
	return r.repo.Find(ctx, ID)
}

func (r *documentRepository) InsertOne(ctx context.Context, m *model.Document) (*model.Document, error) {
	return r.repo.Insert(ctx, m)
}

func (r *documentRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetName("idx_name"),
		},
		{
			Keys:    bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetName("idx_type"),
		},
		{
			Keys:    bson.D{{Key: "contentType", Value: 1}},
			Options: options.Index().SetName("idx_contentType"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "contentType", Value: 1},
			},
			Options: options.Index().SetName("idx_type_contentType"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}
