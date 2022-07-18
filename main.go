package main

import (
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"go.uber.org/zap"

	"github.com/installer/installer/internal/pkg/handlers"
)

func main() {
	pterm.DefaultBigText.WithLetters(putils.LettersFromString("  INSTL")).Render()

	logger, _ := zap.NewProduction(zap.WithCaller(false))
	defer logger.Sync()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(recover.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	app.Get("/", handlers.RedirectToDocs)
	app.Get("/:user/:repo/:os", handlers.Installation)

	pterm.Fatal.PrintOnError(app.Listen(":80"))
}
