package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/antoniofrisenda/template-service/src/internal/api"
	"github.com/antoniofrisenda/template-service/src/internal/assets/dto"
	"github.com/antoniofrisenda/template-service/src/internal/assets/model"
	"github.com/antoniofrisenda/template-service/src/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var testApp *fiber.App

const (
	testBucket    = "test-bucket"
	testDB        = "testDB"
	accessKey     = "test"
	secretKey     = "test"
	awsRegion     = "us-east-1"
	localstackImg = "localstack/localstack:latest"
	mongoImg      = "mongo:latest"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoC, err := mongodb.Run(ctx, mongoImg)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = testcontainers.TerminateContainer(mongoC)
	}()

	mongoURI, err := mongoC.ConnectionString(ctx)
	if err != nil {
		panic(err.Error())
	}

	localstackC, err := localstack.Run(ctx, localstackImg,
		testcontainers.WithEnv(map[string]string{"SERVICES": "s3"}),
	)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = testcontainers.TerminateContainer(localstackC)
	}()

	mappedPort, err := localstackC.MappedPort(ctx, "4566/tcp")
	if err != nil {
		panic(err.Error())
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		panic(err.Error())
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		panic(err.Error())
	}

	s3Endpoint := "http://" + host + ":" + mappedPort.Port()

	if err := createS3Bucket(ctx, s3Endpoint); err != nil {
		panic(err.Error())
	}

	cfg := &config.Config{
		App: config.AppConfig{Port: "0"},
		MongoDB: config.DBConfig{
			URL: mongoURI,
			DB:  testDB,
		},
		AWS: config.AWSConfig{
			Region:            awsRegion,
			AccessKeyID:       accessKey,
			SecretAccessKeyID: secretKey,
			URL:               s3Endpoint,
			S3BucketName:      testBucket,
		},
		Logger: config.LogConfig{
			Format:     "[${time}] ${status} - ${method} ${path} ${latency}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Local",
		},
	}

	testApp, err = api.Init(cfg)
	if err != nil {
		panic(err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

func createS3Bucket(ctx context.Context, endpoint string) error {
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(awsRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(testBucket),
	})
	return err
}

func TestIntegration_PostTemplate_JSON_StaticText_GetTemplate(t *testing.T) {
	doc := dto.InsertDocument{
		Name:        "static-text-doc",
		Summary:     "summary",
		Type:        model.STATIC,
		Source:      model.TEXT,
		ContentType: model.PLAIN_TEXT,
		Body:        &dto.InsertBody{Text: ptr("Hello {{world}}")},
	}
	body, _ := json.Marshal(doc)

	req := httptest.NewRequest(http.MethodPost, "/api/internal/templates/STATIC/TEXT/v1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var created dto.Document
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "static-text-doc", created.Name)
	assert.Equal(t, "Hello {{world}}", created.Body)

	getReq := httptest.NewRequest(http.MethodGet, "/api/internal/templates/STATIC/TEXT/"+created.ID+"/v1", nil)
	getResp, err := testApp.Test(getReq, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	var got dto.Document
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "Hello {{world}}", got.Body)
}

func TestIntegration_GetLatestVariables(t *testing.T) {
	doc := dto.InsertDocument{
		Name:        "template-with-vars",
		Summary:     "summary",
		Type:        model.TEMPLATE,
		Source:      model.TEXT,
		ContentType: model.PLAIN_TEXT,
		Body:        &dto.InsertBody{Text: ptr("Dear {{name}}, your order {{order_id}} is ready.")},
	}
	body, _ := json.Marshal(doc)

	req := httptest.NewRequest(http.MethodPost, "/api/internal/templates/TEMPLATE/TEXT/v1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var created dto.Document
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))

	varsReq := httptest.NewRequest(http.MethodGet, "/api/internal/templates/variables/latest/"+created.ID+"/v1", nil)
	varsResp, err := testApp.Test(varsReq, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer varsResp.Body.Close()

	assert.Equal(t, http.StatusOK, varsResp.StatusCode)
	var varsPayload struct {
		Variables []string `json:"variables"`
	}
	require.NoError(t, json.NewDecoder(varsResp.Body).Decode(&varsPayload))
	assert.Contains(t, varsPayload.Variables, "name")
	assert.Contains(t, varsPayload.Variables, "order_id")
}

func TestIntegration_PostTemplate_MultipartFile_GetTemplate_GetPresigned(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("name", "file-doc")
	_ = w.WriteField("summary", "summary")
	_ = w.WriteField("contentType", "PLAIN_TEXT")

	fw, _ := w.CreateFormFile("file", "hello.txt")
	_, _ = fw.Write([]byte("Hello {{user}} from file"))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/internal/templates/TEMPLATE/FILE/v1", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var created dto.Document
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	assert.NotEmpty(t, created.ID)
	assert.True(t, created.Base64Encoded)

	getReq := httptest.NewRequest(http.MethodGet, "/api/internal/templates/TEMPLATE/FILE/"+created.ID+"/v1", nil)
	getResp, err := testApp.Test(getReq, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	var got dto.Document
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	assert.NotEmpty(t, got.Body)

	presignReq := httptest.NewRequest(http.MethodGet, "/api/internal/templates/url/"+created.ID+"/v1", nil)
	presignResp, err := testApp.Test(presignReq, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer presignResp.Body.Close()

	assert.Equal(t, http.StatusOK, presignResp.StatusCode)
	var presignPayload struct {
		URL string `json:"url"`
	}
	require.NoError(t, json.NewDecoder(presignResp.Body).Decode(&presignPayload))
	assert.NotEmpty(t, presignPayload.URL)
}

func TestIntegration_GetTemplate_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/internal/templates/STATIC/TEXT/000000000000000000000000/v1", nil)
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_GetTemplate_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/internal/templates/STATIC/TEXT/invalid-id/v1", nil)
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_GetPresigned_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/internal/templates/url/invalid/v1", nil)
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestIntegration_PostTemplate_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/internal/templates/STATIC/TEXT/v1", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, fiber.TestConfig{Timeout: 0})
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func ptr(s string) *string { return &s }
