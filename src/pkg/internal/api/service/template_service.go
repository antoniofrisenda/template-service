package service

import (
	"context"
	"fmt"
	"time"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/repository"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service/aws"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/helper"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
	"github.com/gofiber/fiber/v3/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITemplateService interface {
	Find(ctx context.Context, id string) (*dto.TemplatePayload, error)

	SearchTemplateName(ctx context.Context, name string) (*dto.TemplatePayload, error)

	SearchTemplateSummary(ctx context.Context, filter string) (*dto.TemplatePayload, error)
	Create(ctx context.Context, payload dto.TemplatePayload) (*dto.TemplatePayload, error)

	CreateByUploadingFile(ctx context.Context, input []byte, name string, summary string, templateType model.TemplateType, contentType model.ContentType) (*dto.TemplatePayload, error)
	Patch(ctx context.Context, id string, payload dto.TemplatePayload) (*dto.TemplatePayload, error)
	Delete(ctx context.Context, id string) (bool, error)

	UploadBase64(ctx context.Context, key string, data string, contentType string) error
	UploadBytes(ctx context.Context, key string, data []byte, contentType string) error

	DownloadBase64(ctx context.Context, key string) (string, error)
	DownloadBytes(ctx context.Context, key string) ([]byte, error)

	DownloadByPresignedURL(ctx context.Context, key string) (string, error)
}

type TemplateService struct {
	repo     repository.ITemplateRepository
	resolver IResolver
	s3       aws.IS3ClientService
	mapper   helper.ITemplateMapper
}

func (s *TemplateService) Find(ctx context.Context, id string) (*dto.TemplatePayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Find] status=started id=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Find] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	entity, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Find] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	payload, err := s.mapper.ToPayload(*entity)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Find] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Find] status=success id=%s duration=%v", id, time.Since(start))
	return payload, nil
}

func (s *TemplateService) SearchTemplateName(ctx context.Context, name string) (*dto.TemplatePayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.FindByName] status=started name=%s", name)

	entity, err := s.repo.GetByName(ctx, name)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.FindByName] status=failed id=%s error=%v duration=%v", name, err, time.Since(start))
		return nil, err
	}

	payload, err := s.mapper.ToPayload(*entity)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.FindByName] status=failed id=%s error=%v duration=%v", name, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.FindByName] status=success id=%s duration=%v", name, time.Since(start))
	return payload, nil
}

func (s *TemplateService) SearchTemplateSummary(ctx context.Context, filter string) (*dto.TemplatePayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.FindByName] status=started filter=%s", filter)

	entity, err := s.repo.GetBySummary(ctx, filter)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.FindByName] status=failed id=%s error=%v duration=%v", filter, err, time.Since(start))
		return nil, err
	}

	payload, err := s.mapper.ToPayload(*entity)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.FindByName] status=failed id=%s error=%v duration=%v", filter, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.FindByName] status=success id=%s duration=%v", filter, time.Since(start))
	return payload, nil
}

func (s *TemplateService) Create(ctx context.Context, payload dto.TemplatePayload) (*dto.TemplatePayload, error) {
	start := time.Now()
	log.WithContext(ctx).Info("[TemplateService.Create] status=started")

	entity, err := s.mapper.ToTemplate(payload)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] status=failed error=%v duration=%v", err, time.Since(start))
		return nil, err
	}

	if err := helper.ValidateTemplate(*entity); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] status=failed error=%v duration=%v", err, time.Since(start))
		return nil, err
	}

	id, err := s.repo.InsertIntoDB(ctx, entity)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] status=failed error=%v duration=%v", err, time.Since(start))
		return nil, err
	}

	create, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] status=failed id=%s error=%v duration=%v", id.Hex(), err, time.Since(start))
		return nil, err
	}

	result, err := s.mapper.ToPayload(*create)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] status=failed id=%s error=%v duration=%v", id.Hex(), err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Create] status=success id=%s duration=%v", id.Hex(), time.Since(start))
	return result, nil
}

func (s *TemplateService) CreateByUploadingFile(ctx context.Context, input []byte, name string, summary string, templateType model.TemplateType, contentType model.ContentType) (*dto.TemplatePayload, error) {
	fileText, err := s.resolver.ParseFileContent(ctx, input)
	if err != nil {
		return nil, err
	}

	var extractedVars []string
	if contentType == model.HTML || contentType == model.PLAIN_TEXT {
		extractedVars, err = s.resolver.ExtractVariables(ctx, fileText, string(contentType))
		if err != nil {
			return nil, err
		}
	}

	resourcePayload := dto.ResourcePayload{
		Variables: extractedVars,
	}

	key := primitive.NewObjectID().Hex()

	if contentType == model.PLAIN_TEXT {
		resourcePayload.Text = fileText
	} else {
		resourcePayload.URL = fmt.Sprintf("s3://%s/%s", s.s3.GetBucket(), key)
	}

	create, err := s.Create(ctx, dto.TemplatePayload{
		ID:       key,
		Name:     name,
		Summary:  summary,
		Type:     templateType,
		Content:  contentType,
		Resource: resourcePayload,
	})
	if err != nil {
		return nil, err
	}

	if contentType != model.PLAIN_TEXT {
		if err := s.s3.UploadBytes(ctx, key, input, string(contentType)); err != nil {
			return nil, err
		}
	}

	return create, nil
}

func (s *TemplateService) Patch(ctx context.Context, id string, payload dto.TemplatePayload) (*dto.TemplatePayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Patch] status=started id=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	entity, err := s.mapper.ToTemplate(payload)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	if err := helper.ValidateTemplate(*entity); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	update := bson.M{}
	data, err := bson.Marshal(entity)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s marshal error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}
	if err := bson.Unmarshal(data, &update); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s unmarshal error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	if err := s.repo.UpdateByID(ctx, objectID, update); err != nil {
		log.WithContext(ctx).Errorf(
			"[TemplateService.Patch] status=failed id=%s patch error=%v duration=%v",
			id, err, time.Since(start),
		)
		return nil, err
	}

	updated, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s reload error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	result, err := s.mapper.ToPayload(*updated)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] status=failed id=%s mapping response error=%v duration=%v", id, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Patch] status=success id=%s duration=%v", id, time.Since(start))
	return result, nil
}
func (s *TemplateService) Delete(ctx context.Context, id string) (bool, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Delete] status=started id=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Delete] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return false, err
	}

	if err := s.repo.DeleteByID(ctx, objectID); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Delete] status=failed id=%s error=%v duration=%v", id, err, time.Since(start))
		return false, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Delete] status=success id=%s duration=%v", id, time.Since(start))
	return true, nil
}

func (s *TemplateService) UploadBase64(ctx context.Context, key string, data string, contentType string) error {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.UploadBase64] status=started key=%s", key)

	err := s.s3.UploadBase64(ctx, key, data, contentType)

	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.UploadBase64] status=failed key=%s error=%v duration=%v", key, err, time.Since(start))
		return err
	}

	log.WithContext(ctx).Infof("[TemplateService.UploadBase64] status=success key=%s duration=%v", key, time.Since(start))
	return nil
}

func (s *TemplateService) UploadBytes(ctx context.Context, key string, data []byte, contentType string) error {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.UploadBytes] status=started key=%s", key)

	err := s.s3.UploadBytes(ctx, key, data, contentType)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.UploadBytes] status=failed key=%s error=%v duration=%v", key, err, time.Since(start))
		return err
	}

	log.WithContext(ctx).Infof("[TemplateService.UploadBytes] status=success key=%s duration=%v", key, time.Since(start))
	return nil
}

func (s *TemplateService) DownloadBase64(ctx context.Context, key string) (string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.DownloadBase64] status=started key=%s", key)

	data, err := s.s3.DownloadBase64(ctx, key)
	if err != nil {
		log.WithContext(ctx).Errorf(
			"[TemplateService.DownloadBase64] status=failed key=%s error=%v duration=%v",
			key, err, time.Since(start),
		)
		return "", err
	}

	log.WithContext(ctx).Infof("[TemplateService.DownloadBase64] status=success key=%s duration=%v", key, time.Since(start))
	return data, nil
}

func (s *TemplateService) DownloadBytes(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.DownloadBytes] status=started key=%s", key)

	data, err := s.s3.DownloadBytes(ctx, key)
	if err != nil {
		log.WithContext(ctx).Errorf(
			"[TemplateService.DownloadBytes] status=failed key=%s error=%v duration=%v",
			key, err, time.Since(start),
		)
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.DownloadBytes] status=success key=%s duration=%v", key, time.Since(start))
	return data, nil
}

func (s *TemplateService) DownloadByPresignedURL(ctx context.Context, key string) (string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.DownloadByPresignedURL] status=started id=%s", key)

	url, err := s.s3.GetPresignedURL(ctx, key, 10*time.Minute)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.DownloadByPresignedURL] status=failed id=%s error=%v duration=%v", key, err, time.Since(start))
		return "", err
	}

	log.WithContext(ctx).Infof("[TemplateService.Find] status=success id=%s duration=%v", key, time.Since(start))
	return url, err
}

func NewTemplateService(s3 aws.IS3ClientService, mapper helper.ITemplateMapper, repo repository.ITemplateRepository, resolver IResolver) ITemplateService {
	return &TemplateService{
		repo:     repo,
		s3:       s3,
		mapper:   mapper,
		resolver: resolver,
	}
}
