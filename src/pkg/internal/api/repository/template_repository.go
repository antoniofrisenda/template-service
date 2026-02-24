package repository

import (
	"context"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITemplateRepository interface {
	Insert(ctx context.Context, m *model.Template) (primitive.ObjectID, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Template, error)
	FindByName(ctx context.Context, s string) (*model.Template, error)
	FindBySummary(ctx context.Context, s string) (*model.Template, error)
	Patch(ctx context.Context, id primitive.ObjectID, update bson.M) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type TemplateRepository struct {
	repo *CRUDRepository[model.Template]
}

func NewTemplateRepository(collection *mongo.Collection) ITemplateRepository {
	return &TemplateRepository{repo: NewGenericRepository[model.Template](collection)}
}

func (r *TemplateRepository) Insert(ctx context.Context, m *model.Template) (primitive.ObjectID, error) {
	return r.repo.Create(ctx, m)
}

func (r *TemplateRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Template, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *TemplateRepository) FindByName(ctx context.Context, s string) (*model.Template, error) {
	var t model.Template
	if err := r.repo.collection.FindOne(ctx, bson.M{"name": s}).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) FindBySummary(ctx context.Context, s string) (*model.Template, error) {
	var t model.Template
	filter := bson.M{"summary": bson.M{"$regex": s, "$options": "i"}}
	if err := r.repo.collection.FindOne(ctx, filter).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) Patch(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	return r.repo.UpdateByID(ctx, id, update)
}

func (r *TemplateRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	return r.repo.DeleteByID(ctx, id)
}
