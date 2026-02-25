package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"strings"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type IS3ClientService interface {
	GetBucket(ctx context.Context) error
	UploadBase64(ctx context.Context, key string, data string, contentType string) error
	UploadBytes(ctx context.Context, key string, data []byte, contentType string) error
	DownloadBase64(ctx context.Context, key string) (string, error)
	DownloadBytes(ctx context.Context, key string) ([]byte, error)
	GetPresignedURL(ctx context.Context, key string, lifetime time.Duration) (string, error)
}

type S3ClientService struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	region    string
}

func NewS3ClientService(ctx context.Context, region string, bucket string, access_key string, secret_key string) (IS3ClientService, error) {
	var cfgOpts []func(*awsConfig.LoadOptions) error
	cfgOpts = append(cfgOpts, awsConfig.WithRegion(region))
	if k := access_key; k != "" {
		cfgOpts = append(cfgOpts, awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(k, secret_key, "")))
	}

	cfg, err := awsConfig.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3ClientService{
		client:    client,
		presigner: s3.NewPresignClient(client),
		bucket:    bucket,
		region:    region,
	}, nil
}

func (s *S3ClientService) GetBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &s.bucket})
	if err == nil {
		return nil
	}

	create := &s3.CreateBucketInput{Bucket: &s.bucket}
	create.CreateBucketConfiguration = &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraint(s.region)}

	_, err = s.client.CreateBucket(ctx, create)
	if err != nil {
		if strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			return nil
		}
		return err
	}

	return s3.NewBucketExistsWaiter(s.client).Wait(ctx, &s3.HeadBucketInput{Bucket: &s.bucket}, 30*time.Second)
}

func (s *S3ClientService) UploadBase64(ctx context.Context, key string, data string, contentType string) error {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return s.UploadBytes(ctx, key, decoded, contentType)
}

func (s *S3ClientService) UploadBytes(ctx context.Context, key string, stream []byte, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.bucket,
		Key:         &key,
		Body:        bytes.NewReader(stream),
		ContentType: &contentType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3ClientService) DownloadBase64(ctx context.Context, key string) (string, error) {
	data, err := s.DownloadBytes(ctx, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *S3ClientService) DownloadBytes(ctx context.Context, key string) ([]byte, error) {
	stream, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		return nil, err
	}
	defer stream.Body.Close()

	bytes, err := io.ReadAll(stream.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (s *S3ClientService) GetPresignedURL(ctx context.Context, key string, lifetime time.Duration) (string, error) {
	url, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key}, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", err
	}
	return url.URL, nil
}
