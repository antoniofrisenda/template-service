package cmd

import (
	"log"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := middleware.Init()
	log.Fatal(app.Listen(":3000", fiber.ListenConfig{
		DisableStartupMessage: true,
	}))
}
