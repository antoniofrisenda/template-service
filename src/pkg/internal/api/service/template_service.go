package service

import (
	"context"
	"time"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/repository"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service/aws"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/helper"
	"github.com/gofiber/fiber/v3/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITemplateService interface {
	Find(ctx context.Context, id string) (*dto.TemplatetPayload, error)
	Create(ctx context.Context, dto dto.TemplatetPayload) (*dto.TemplatetPayload, error)

	Patch(ctx context.Context, id string, dto dto.TemplatetPayload) (*dto.TemplatetPayload, error)

	Delete(ctx context.Context, id string) (bool, error)

	UploadBase64(ctx context.Context, key string, data string, contentType string) error

	UploadBytes(ctx context.Context, key string, stream []byte, contentType string) error

	DownloadBase64(ctx context.Context, key string) (string, error)

	DownloadBytes(ctx context.Context, key string) ([]byte, error)

	GetPresignedURL(ctx context.Context, key string) (string, error)
}

type TemplateService struct {
	repo repository.ITemplateRepository
	s3   aws.IS3ClientService
}

func NewTemplateService(repo repository.ITemplateRepository, s3 aws.IS3ClientService) ITemplateService {
	return &TemplateService{
		repo: repo,
		s3:   s3,
	}
}

//CRUD

func (s *TemplateService) Find(ctx context.Context, id string) (*dto.TemplatetPayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Find] status=started target=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Find] step=ObjectIDFromHex(id) status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	find, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Find] step=repo.FindByID(ctx, objectID) status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Find] status=success step=done target_id=%s duration=%v", id, time.Since(start))
	return helper.ToPayload(*find), nil

}

func (s *TemplateService) Create(ctx context.Context, dto dto.TemplatetPayload) (*dto.TemplatetPayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Create] status=started")

	t := helper.ToTemplate(dto)
	if err := helper.ValidateTemplate(t); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] step=ValidateTemplate status=failed error=%q duration=%v", err, time.Since(start))
		return nil, err
	}

	id, err := s.repo.Insert(ctx, &t)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] step=repo.Insert status=failed error=%q duration=%v", err, time.Since(start))
		return nil, err
	}

	created, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Create] step=repo.FindByID status=failed target_id=%s error=%q duration=%v", id.Hex(), err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Create] status=success target_id=%s duration=%v", id.Hex(), time.Since(start))
	return helper.ToPayload(*created), nil
}

func (s *TemplateService) Delete(ctx context.Context, id string) (bool, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Delete] status=started target=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Delete] step=ObjectIDFromHex(id) status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return false, err
	}

	if err := s.repo.DeleteByID(ctx, objectID); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Delete] step=repo.DeleteByID status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return false, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Delete] status=success target_id=%s duration=%v", id, time.Since(start))
	return true, nil
}

func (s *TemplateService) Patch(ctx context.Context, id string, dto dto.TemplatetPayload) (*dto.TemplatetPayload, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.Patch] status=started target=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] step=ObjectIDFromHex(id) status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	t := helper.ToTemplate(dto)
	if err := helper.ValidateTemplate(t); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] step=ValidateTemplate status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	update := bson.M{
		"name":    t.Name,
		"summary": t.Summary,
		"type":    t.Type,
		"content": t.Content,
		"resource": bson.M{
			"url":       t.Resource.URL,
			"text":      t.Resource.Text,
			"variables": t.Resource.Variables,
		},
	}

	if err := s.repo.Patch(ctx, objectID, update); err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] step=repo.Patch status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	updated, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.Patch] step=repo.FindByID status=failed target_id=%s error=%q duration=%v", id, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.Patch] status=success target_id=%s duration=%v", id, time.Since(start))
	return helper.ToPayload(*updated), nil
}

//S3

func (s *TemplateService) UploadBase64(ctx context.Context, key string, data string, contentType string) error {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.UploadBase64] status=started target=%s", key)

	err := s.s3.UploadBase64(ctx, key, data, contentType)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.UploadBase64] step=s3.UploadBase64 status=failed target=%s error=%q duration=%v", key, err, time.Since(start))
		return err
	}

	log.WithContext(ctx).Infof("[TemplateService.UploadBase64] status=success target=%s duration=%v", key, time.Since(start))
	return nil
}

func (s *TemplateService) UploadBytes(ctx context.Context, key string, stream []byte, contentType string) error {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.UploadBytes] status=started target=%s", key)

	err := s.s3.UploadBytes(ctx, key, stream, contentType)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.UploadBytes] step=s3.UploadBytes status=failed target=%s error=%q duration=%v", key, err, time.Since(start))
		return err
	}

	log.WithContext(ctx).Infof("[TemplateService.UploadBytes] status=success target=%s duration=%v", key, time.Since(start))
	return nil
}

func (s *TemplateService) DownloadBase64(ctx context.Context, key string) (string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.DownloadBase64] status=started target=%s", key)

	data, err := s.s3.DownloadBase64(ctx, key)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.DownloadBase64] step=s3.DownloadBase64 status=failed target=%s error=%q duration=%v", key, err, time.Since(start))
		return "", err
	}

	log.WithContext(ctx).Infof("[TemplateService.DownloadBase64] status=success target=%s duration=%v", key, time.Since(start))
	return data, nil
}

func (s *TemplateService) DownloadBytes(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.DownloadBytes] status=started target=%s", key)

	data, err := s.s3.DownloadBytes(ctx, key)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.DownloadBytes] step=s3.DownloadBytes status=failed target=%s error=%q duration=%v", key, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[TemplateService.DownloadBytes] status=success target=%s duration=%v", key, time.Since(start))
	return data, nil
}

func (s *TemplateService) GetPresignedURL(ctx context.Context, key string) (string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[TemplateService.GetPresignedURL] status=started target=%s", key)

	url, err := s.s3.GetPresignedURL(ctx, key, 15*time.Minute)
	if err != nil {
		log.WithContext(ctx).Errorf("[TemplateService.GetPresignedURL] step=s3.GetPresignedURL status=failed target=%s error=%q duration=%v", key, err, time.Since(start))
		return "", err
	}

	log.WithContext(ctx).Infof("[TemplateService.GetPresignedURL] status=success target=%s duration=%v", key, time.Since(start))
	return url, nil
}
