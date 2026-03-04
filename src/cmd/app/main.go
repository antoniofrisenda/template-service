package main

import (
	"log"

	"github.com/antoniofrisenda/template-service/src/internal/api"
	"github.com/antoniofrisenda/template-service/src/internal/config"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app, err := api.Init(cfg)
	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":"+cfg.App.Port, fiber.ListenConfig{
		DisableStartupMessage: true,
	}))
}
