package main

import (
	"log"

	"github.com/antoniofrisenda/template-service/src/internal/api"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app, err := api.Init()
	if err != nil {
		panic(err)
	}

	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		DisableStartupMessage: true,
	}))
}
