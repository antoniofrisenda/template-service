package aws

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	GetBucket() string
	EnsureBucketExists(ctx context.Context) error
	Upload(ctx context.Context, key string, body io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	DownloadWithPresignedURL(ctx context.Context, key string, lifetime time.Duration) (string, error)
}

type s3Client struct {
	Bucket        string
	Region        string
	S3            *s3.Client
	PresignClient *s3.PresignClient
}

func (s *s3Client) GetBucket() string {
	return s.Bucket
}

func (s *s3Client) EnsureBucketExists(ctx context.Context) error {
	_, err := s.S3.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket),
	})
	return err
}

func (s *s3Client) Upload(ctx context.Context, key string, body io.Reader) error {
	_, err := s.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	return err
}

func (s *s3Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	body, err := s.S3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return body.Body, nil
}

func (s *s3Client) DownloadWithPresignedURL(ctx context.Context, key string, lifetime time.Duration) (string, error) {
	url, err := s.PresignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(lifetime),
	)

	if err != nil {
		return "", err
	}

	return url.URL, nil
}

func NewS3ClientService(ctx context.Context, region, accessKey, secretKey, endpoint, bucket string) (S3Client, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return &s3Client{
		Bucket:        bucket,
		Region:        region,
		S3:            client,
		PresignClient: s3.NewPresignClient(client),
	}, nil
}
