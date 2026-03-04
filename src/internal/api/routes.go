package api

import (
	"context"
	"encoding/json"

	AWS "github.com/antoniofrisenda/template-service/src/clients/aws"
	MONGO "github.com/antoniofrisenda/template-service/src/clients/mongo"
	"github.com/antoniofrisenda/template-service/src/internal/api/router"
	"github.com/antoniofrisenda/template-service/src/internal/assets/helpers"
	"github.com/antoniofrisenda/template-service/src/internal/config"
	"github.com/antoniofrisenda/template-service/src/internal/repository"
	"github.com/antoniofrisenda/template-service/src/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func Init(cfg *config.Config) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(cfg.Logger.NewFiberLogger())

	app.Use(requestid.New())

	log.Info("Init app...")

	ctx := context.Background()

	err := RegisterInternalRoute(ctx, cfg, app)
	if err != nil {
		return nil, err
	}

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "document-service",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	log.Info("Init app OK!")

	return app, nil
}

func RegisterInternalRoute(ctx context.Context, cfg *config.Config, app *fiber.App) error {
	route := app.Group("/api/internal/templates")

	mongoClient, err := MONGO.NewMongoClient(
		ctx,
		cfg.MongoDB.URL,
		cfg.MongoDB.DB,
	)
	if err != nil {
		panic(err)
	}

	repo := repository.NewDocumentRepository(mongoClient.GetDB().Collection("templates"))

	s3, err := AWS.NewS3ClientService(
		ctx,
		cfg.AWS.Region,
		cfg.AWS.AccessKeyID,
		cfg.AWS.SecretAccessKeyID,
		cfg.AWS.URL,
		cfg.AWS.S3BucketName,
	)
	if err != nil {
		panic(err)
	}

	if err := s3.EnsureBucketExists(ctx); err != nil {
		panic(err)
	} else {
		log.Info("S3 bucket OK!")
	}

	mapper := helpers.NewDocumentMapper()

	service := service.NewDocumentService(repo, mapper, s3)

	controller := router.NewDocumentController(service)

	route.Get("/url/:ID/v1", controller.GetPresigned)
	route.Get("/variables/latest/:ID/v1", controller.GetLatestVariables)
	route.Get("/:DocumentType/:SourceType/:ID/v1", controller.GetTemplate)
	route.Post("/:DocumentType/:SourceType/v1", controller.PostTemplate)

	return nil
}
