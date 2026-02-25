package middleware

import (
	"context"
	"encoding/json"
	"os"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/repository"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/service"
	awsService "github.com/antoniofrisenda/template-service/src/pkg/internal/api/service/aws"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/api/v1/router"
	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/factory/helper"
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

	templateMapper := helper.NewTemplateMapper(helper.NewRegistry())

	controller := router.NewTemplateController(
		service.NewTemplateService(
			s3Client,
			templateMapper,
			repository.NewTemplateRepository(client.GetConnection().Collection("templates")),
			service.NewResolver(),
		),
	)

	route := app.Group("/api/internal/templates")
	route.Get("/:id/v1", controller.FindTemplate)
	route.Get("/name/v1", controller.FindTemplateByName)
	route.Get("/summary/v1", controller.FindTemplateBySummary)
	route.Get("/s3/:id/download/bytes/v1", controller.DownloadFileByBytes)
	route.Get("/:id/download/base64/v1", controller.DownloadFileByBase64)
	route.Get("/:id/download/presigned-url/v1", controller.DownloadFileByPresignedURL)
	route.Post("/v1", controller.CreateTemplate)
	route.Post("/s3/upload/v1", controller.UploadFile)
	route.Patch("/:id/v1", controller.PatchTemplate)
	route.Delete("/:id/v1", controller.DeleteTemplate)
}
