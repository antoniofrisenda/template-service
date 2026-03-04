package repository

import (
	"context"

	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DocumentRepository interface {
	FindOne(ctx context.Context, ID primitive.ObjectID) (*model.Document, error)
	InsertOne(ctx context.Context, m *model.Document) (*model.Document, error)
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
