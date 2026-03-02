package config

import (
	"fmt"
	"os"
)

type Config struct {
	App     ServerConfig
	MongoDB DatabaseConfig
	AWS     AWSConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
	DB  string
}

type AWSConfig struct {
	Region            string
	AccessKeyID       string
	SecretAccessKeyID string
	URL               string
	S3BucketName      string
}

func Load() (*Config, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}
	url, err := getEnv("MONGO_URL")
	if err != nil {
		return nil, err
	}
	db, err := getEnv("DB")
	if err != nil {
		return nil, err
	}
	awsRegion, err := getEnv("AWS_DEFAULT_REGION")
	if err != nil {
		return nil, err
	}
	awsAccessKey, err := getEnv("AWS_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}
	awsSecretKey, err := getEnv("AWS_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}
	awsEndpoint, _ := getEnv("AWS_ENDPOINT_URL")
	awsBucket, err := getEnv("AWS_S3_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: ServerConfig{
			Port: port,
		},
		MongoDB: DatabaseConfig{
			URL: url,
			DB:  db,
		},
		AWS: AWSConfig{
			Region:            awsRegion,
			AccessKeyID:       awsAccessKey,
			SecretAccessKeyID: awsSecretKey,
			URL:               awsEndpoint,
			S3BucketName:      awsBucket,
		},
	}

	return cfg, nil
}

func getEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("%s is required", key)
}
