package config

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type Config struct {
	App     AppConfig
	MongoDB DBConfig
	AWS     AWSConfig
	Logger  LogConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
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

type LogConfig struct {
	Format     string
	TimeFormat string
	TimeZone   string
}

func (l LogConfig) NewFiberLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     l.Format,
		TimeFormat: l.TimeFormat,
		TimeZone:   l.TimeZone,
	})
}

func Load() (*Config, error) {
	port, err := Get("PORT", "")
	if err != nil {
		return nil, err
	}

	url, err := Get("MONGO_URL", "")
	if err != nil {
		return nil, err
	}

	db, err := Get("DB", "")
	if err != nil {
		return nil, err
	}

	awsRegion, err := Get("AWS_DEFAULT_REGION", "")
	if err != nil {
		return nil, err
	}

	awsAccessKey, err := Get("AWS_ACCESS_KEY_ID", "")
	if err != nil {
		return nil, err
	}

	awsSecretKey, err := Get("AWS_SECRET_ACCESS_KEY", "")
	if err != nil {
		return nil, err
	}

	awsEndpoint, err := Get("AWS_ENDPOINT_URL", "")
	if err != nil {
		return nil, err
	}

	awsBucket, err := Get("AWS_S3_BUCKET_NAME", "")
	if err != nil {
		return nil, err
	}

	loggerFormat, err := Get("LOGGER_FORMAT", "[${time}] ${status} - ${method} ${path} ${latency}\n")
	if err != nil {
		return nil, err
	}

	loggerTimeFormat, err := Get("LOGGER_TIME_FORMAT", "2006-01-02 15:04:05")
	if err != nil {
		return nil, err
	}

	loggerTimeZone, err := Get("LOGGER_TIME_ZONE", "Local")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{Port: port},
		MongoDB: DBConfig{
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
		Logger: LogConfig{
			Format:     loggerFormat,
			TimeFormat: loggerTimeFormat,
			TimeZone:   loggerTimeZone,
		},
	}

	return cfg, nil
}

func Get(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	if fallback != "" {
		return fallback, nil
	}

	return "", fmt.Errorf("%s is required", key)
}
