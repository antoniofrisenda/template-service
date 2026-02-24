package middleware

import (
	"context"
	"encoding/json"
	"os"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/repository"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service"
	awsService "github.com/antoniofrisenda/template-service/src/pkg/internal/api/service/aws"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/v1/router"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/dto"
	mongoClient "github.com/antoniofrisenda/template-service/src/pkg/internal/assets/mongo"
	LogConfig "github.com/antoniofrisenda/template-service/src/pkg/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(requestid.New())

	app.Use(logger.New(LogConfig.Format))

	log.Info("Init OK!")

	RegisterExternalRoute(app)
	RegisterInternalRoute(app, context.Background())

	return app
}

func RegisterExternalRoute(app *fiber.App) {
	route := app.Group("/api")

	route.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"detail": "found",
		})
	})

	route.Get("/healthz", func(ctx fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})

}

func RegisterInternalRoute(app *fiber.App, ctx context.Context) {

	client := mongoClient.NewMongoClient(os.Getenv("MONGO_URI"), os.Getenv("MONGO_DB_NAME"))
	if err := client.Connect(); err != nil {
		panic(err)
	}

	s3Client, err := awsService.NewS3ClientService(
		ctx,
		os.Getenv("AWS_DEFAULT_REGION"),
		os.Getenv("AWS_BUCKET_NAME"),
		os.Getenv("AWS_KEY"),
		os.Getenv("AWS_SECRET_KEY"),
	)
	if err != nil {
		panic(err)
	}
	if err := s3Client.GetBucket(ctx); err != nil {
		panic(err)
	}

	controller := router.NewTemplateController(service.NewTemplateService(repository.NewTemplateRepository(client.GetConnection().Collection("templates")), s3Client))

	route := app.Group("/api/internal/templates")
	route.Post("/v1", controller.Create)
	route.Post("/s3/upload/v1", func(c fiber.Ctx) error {
		var b64 dto.S3UploadBase64Payload
		if err := c.Bind().Body(&b64); err == nil && b64.FileName != "" && b64.Base64Data != "" {
			return controller.Upload(c, b64.Base64Data, b64.ContentType)
		}

		var bytesPayload dto.S3UploadBytesPayload
		if err := c.Bind().Body(&bytesPayload); err == nil && bytesPayload.FileName != "" && len(bytesPayload.Bytes) > 0 {
			return controller.Upload(c, bytesPayload.Bytes, bytesPayload.ContentType)
		}

		return fiber.NewError(fiber.StatusBadRequest, "invalid upload payload")
	})

	route.Get("/:id/v1", controller.Find)
	route.Get("/s3/download/:key/v1", controller.Download)
	route.Get("/s3/presigned/:key/v1", controller.PresignedURL)

	route.Patch("/:id/v1", controller.Patch)
	route.Delete("/:id/v1", controller.Delete)
}
