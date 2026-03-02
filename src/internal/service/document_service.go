package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"

	"github.com/antoniofrisenda/template-service/src/clients/aws"
	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/helpers"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
	"github.com/antoniofrisenda/template-service/src/internal/repository"
	"github.com/unidoc/unipdf/v3/extractor"
	unipdfmodel "github.com/unidoc/unipdf/v3/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var regex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}}`)

type DocumentService interface {
	ExtractVariables(ctx context.Context, ID string) ([]string, error)
	FindTemplate(ctx context.Context, ID string) (*dto.Document, error)

	FindTemplateWithPresignedURL(ctx context.Context, ID string) (string, error)
	InsertTemplate(ctx context.Context, d *dto.InsertDocument, file *multipart.FileHeader) (*dto.Document, error)
}

type documentService struct {
	repo   repository.DocumentRepository
	mapper helpers.DocumentMapper
	s3     aws.S3Client
}

func (d *documentService) ExtractVariables(ctx context.Context, ID string) ([]string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[DocumentService.ExtractVariables] status=started target=%s", ID)

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.ExtractVariables] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, err
	}

	doc, err := d.repo.FindOne(ctx, objID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.ExtractVariables] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, err
	}

	extracted, err := d.extractVariables(ctx, doc)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.ExtractVariables] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, err
	}

	log.WithContext(ctx).Infof("[DocumentService.ExtractVariables] status=success target=%s duration=%s", ID, time.Since(start))
	return extracted, nil
}

func (d *documentService) FindTemplate(ctx context.Context, ID string) (*dto.Document, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[DocumentService.FindTemplate] status=started target=%s", ID)

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, fmt.Errorf("invalid object id: %w", err)
	}

	doc, err := d.repo.FindOne(ctx, objID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, fmt.Errorf("document not found: %w", err)
	}

	var result *dto.Document

	switch doc.Source {
	case model.TEXT:
		result, err = d.mapper.ToDTO(doc)
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
			return nil, err
		}

		log.WithContext(ctx).Infof("[DocumentService.FindTemplate] status=success target=%s duration=%s", ID, time.Since(start))
		return result, nil

	case model.FILE:
		reader, err := d.s3.Download(ctx, *doc.Body.URL)
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
			return nil, fmt.Errorf("failed to download file: %w", err)
		}
		defer reader.Close()

		content, err := io.ReadAll(reader)
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(content)
		doc.Body.URL = &encoded
		doc.Body.Text = nil

		result, err = d.mapper.ToDTO(doc)
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
			return nil, err
		}

		log.WithContext(ctx).Infof("[DocumentService.FindTemplate] status=success target=%s duration=%s", ID, time.Since(start))
		return result, nil

	default:
		err := fmt.Errorf("unsupported source type: %s", doc.Source)
		log.WithContext(ctx).Errorf("[DocumentService.FindTemplate] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return nil, err
	}
}

func (d *documentService) FindTemplateWithPresignedURL(ctx context.Context, ID string) (string, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[DocumentService.FindWithPresignedURL] status=started target=%s", ID)

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.FindWithPresignedURL] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return "", fmt.Errorf("invalid object id: %w", err)
	}

	doc, err := d.repo.FindOne(ctx, objID)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.FindWithPresignedURL] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return "", fmt.Errorf("document not found: %w", err)
	}

	if doc.Source != model.FILE || doc.Body.URL == nil {
		err := errors.New("document is not a file")
		log.WithContext(ctx).Errorf("[DocumentService.FindWithPresignedURL] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return "", err
	}

	url, err := d.s3.DownloadWithPresignedURL(ctx, *doc.Body.URL, 15*time.Minute)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.FindWithPresignedURL] status=failure target=%s error=%v duration=%s", ID, err, time.Since(start))
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.WithContext(ctx).Infof("[DocumentService.FindWithPresignedURL] status=success target=%s duration=%s", ID, time.Since(start))
	return url, nil
}

func (d *documentService) InsertTemplate(ctx context.Context, payload *dto.InsertDocument, file *multipart.FileHeader) (*dto.Document, error) {
	start := time.Now()
	log.WithContext(ctx).Infof("[DocumentService.InsertTemplate] status=started")

	doc, err := d.mapper.ToModel(payload)
	if err != nil || doc == nil {
		log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure error=%v duration=%s", err, time.Since(start))
		return nil, fmt.Errorf("failed to map payload to model: %w", err)
	}

	if doc.ID.IsZero() {
		doc.ID = primitive.NewObjectID()
	}

	if doc.Type == model.TEMPLATE && doc.Source == model.TEXT && doc.Body.Text != nil {
		doc.Body.Variables, err = d.extractVariables(ctx, doc)
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure extracting variables error=%v duration=%s", err, time.Since(start))
			return nil, fmt.Errorf("failed to extract variables: %w", err)
		}
	}

	if doc.Source == model.FILE {
		if file == nil {
			log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure error=file is nil for FILE source")
			return nil, errors.New("file is required for document of type FILE")
		}

		src, err := file.Open()
		if err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure error=%v duration=%s", err, time.Since(start))
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer src.Close()

		key := fmt.Sprintf("s3://%s/documents/%s/%s", d.s3.GetBucket(), doc.ID.Hex(), file.Filename)

		if err := d.s3.Upload(ctx, key, src); err != nil {
			log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure uploading to S3 error=%v duration=%s", err, time.Since(start))
			return nil, fmt.Errorf("failed to upload file to S3: %w", err)
		}

		doc.Body.URL = &key

		if doc.Type == model.TEMPLATE {
			doc.Body.Variables, err = d.extractVariables(ctx, doc)
			if err != nil {
				log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure extracting variables error=%v duration=%s", err, time.Since(start))
				return nil, fmt.Errorf("failed to extract variables: %w", err)
			}
		}
	}

	inserted, err := d.repo.InsertOne(ctx, doc)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure inserting to DB error=%v duration=%s", err, time.Since(start))
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	result, err := d.mapper.ToDTO(inserted)
	if err != nil {
		log.WithContext(ctx).Errorf("[DocumentService.InsertTemplate] status=failure converting to DTO error=%v duration=%s", err, time.Since(start))
		return nil, fmt.Errorf("failed to convert to DTO: %w", err)
	}

	log.WithContext(ctx).Infof("[DocumentService.InsertTemplate] status=success target=%s duration=%s", inserted.ID.Hex(), time.Since(start))
	return result, nil
}

func NewDocumentService(repo repository.DocumentRepository, mapper helpers.DocumentMapper, s3 aws.S3Client) DocumentService {
	return &documentService{
		repo:   repo,
		mapper: mapper,
		s3:     s3,
	}
}

func (d *documentService) extractVariables(ctx context.Context, doc *model.Document) ([]string, error) {
	var content string

	switch doc.Source {
	case model.TEXT:
		if doc.Body.Text == nil {
			return nil, fmt.Errorf("text body is nil")
		}
		content = *doc.Body.Text

	case model.FILE:
		if doc.Body.URL == nil {
			return nil, fmt.Errorf("file URL is nil")
		}
		reader, err := d.s3.Download(ctx, *doc.Body.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to download file: %w", err)
		}
		defer reader.Close()

		Bytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		if doc.ContentType == model.PDF {
			pdf, err := unipdfmodel.NewPdfReader(bytes.NewReader(Bytes))
			if err != nil {
				return nil, fmt.Errorf("failed to parse pdf: %w", err)
			}

			pages, err := pdf.GetNumPages()
			if err != nil {
				return nil, fmt.Errorf("failed to read pdf page count: %w", err)
			}

			var textBuilder strings.Builder
			for i := 1; i <= pages; i++ {
				page, err := pdf.GetPage(i)
				if err != nil {
					return nil, fmt.Errorf("failed to get pdf page %d: %w", i, err)
				}

				Exctractor, err := extractor.New(page)
				if err != nil {
					return nil, fmt.Errorf("failed to init extractor for pdf page %d: %w", i, err)
				}

				Text, err := Exctractor.ExtractText()
				if err != nil {
					return nil, fmt.Errorf("failed to extract text from pdf page %d: %w", i, err)
				}

				_, err = textBuilder.WriteString(Text)
				if err != nil {
					return nil, err
				}
			}

			content = textBuilder.String()
		} else {
			content = string(Bytes)
		}
	default:
		return nil, fmt.Errorf("unsupported source type: %s", doc.Source)
	}

	matched := make(map[string]struct{})
	for _, m := range regex.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 {
			matched[m[1]] = struct{}{}
		}
	}

	variables := make([]string, 0, len(matched))
	for k := range matched {
		variables = append(variables, k)
	}

	return variables, nil
}
