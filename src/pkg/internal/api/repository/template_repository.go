package repository

import (
	"context"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITemplateRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Template, error)
	GetByName(ctx context.Context, s string) (*model.Template, error)
	GetBySummary(ctx context.Context, s string) (*model.Template, error)
	InsertIntoDB(ctx context.Context, m *model.Template) (primitive.ObjectID, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type TemplateRepository struct {
	repo *CRUDRepository[model.Template]
}

func NewTemplateRepository(collection *mongo.Collection) ITemplateRepository {
	return &TemplateRepository{repo: NewGenericRepository[model.Template](collection)}
}

func (r *TemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Template, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *TemplateRepository) GetByName(ctx context.Context, s string) (*model.Template, error) {
	var t model.Template
	if err := r.repo.collection.FindOne(ctx, bson.M{"name": s}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) GetBySummary(ctx context.Context, f string) (*model.Template, error) {
	var t model.Template
	filter := bson.M{"summary": bson.M{"$regex": f, "$options": "i"}}
	if err := r.repo.collection.FindOne(ctx, filter).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) InsertIntoDB(ctx context.Context, m *model.Template) (primitive.ObjectID, error) {
	return r.repo.CreateEntity(ctx, m)
}

func (r *TemplateRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	return r.repo.UpdateByID(ctx, id, update)
}

func (r *TemplateRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	return r.repo.DeleteByID(ctx, id)
}
